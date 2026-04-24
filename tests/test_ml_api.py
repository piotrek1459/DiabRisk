from pathlib import Path
import tempfile
import unittest

import numpy as np
from fastapi import HTTPException

from src.FastAPI import ml_api


FEATURE_NAMES = ["HighBP", "BMI", "MentHlth", "PhysHlth", "Age"]


class FakeModel:
    classes_ = np.array([0, 1])

    def __init__(self, at_risk_probability):
        self.at_risk_probability = at_risk_probability
        self.seen_X = None

    def predict_proba(self, X):
        self.seen_X = X.copy()
        return np.array(
            [[1 - self.at_risk_probability, self.at_risk_probability]],
            dtype=float,
        )


class MLApiTests(unittest.TestCase):
    def setUp(self):
        self.original_artifact = ml_api._artifact
        self.original_model_path = ml_api.MODEL_PATH
        self.original_candidates = ml_api.CANDIDATES
        ml_api._artifact = None
        ml_api.MODEL_PATH = None

    def tearDown(self):
        ml_api._artifact = self.original_artifact
        ml_api.MODEL_PATH = self.original_model_path
        ml_api.CANDIDATES = self.original_candidates

    def set_fake_artifact(self, model, model2=None):
        ml_api._artifact = {
            "model1": model,
            "model2": model2,
            "feature_names": FEATURE_NAMES,
            "model_features": "processed",
            "api_input_features": "raw",
        }

    def valid_features(self):
        return {
            "HighBP": 1,
            "BMI": 30.0,
            "MentHlth": 5,
            "PhysHlth": 8,
            "Age": 7,
        }

    def test_build_X_preserves_feature_order_and_casts_to_float(self):
        X = ml_api.build_X({"BMI": "30.5", "Age": 7}, ["BMI", "Age"])

        self.assertEqual(X.shape, (1, 2))
        np.testing.assert_allclose(X, np.array([[30.5, 7.0]]))

    def test_build_X_rejects_missing_features(self):
        with self.assertRaises(HTTPException) as ctx:
            ml_api.build_X({"BMI": 30}, ["BMI", "Age", "Income"])

        self.assertEqual(ctx.exception.status_code, 400)
        self.assertEqual(
            ctx.exception.detail,
            {"error": "Missing features", "missing": ["Age", "Income"]},
        )

    def test_build_X_rejects_non_numeric_values(self):
        with self.assertRaises(HTTPException) as ctx:
            ml_api.build_X({"BMI": "not-a-number"}, ["BMI"])

        self.assertEqual(ctx.exception.status_code, 400)
        self.assertEqual(ctx.exception.detail, {"error": "Invalid feature value type"})

    def test_apply_training_scaling_scales_only_training_scaled_columns(self):
        X = np.array([[1.0, 30.0, 5.0, 8.0, 7.0]])

        scaled = ml_api.apply_training_scaling(X, FEATURE_NAMES)

        expected = X.copy()
        for feature in ("BMI", "MentHlth", "PhysHlth"):
            idx = FEATURE_NAMES.index(feature)
            expected[0, idx] = (
                expected[0, idx] - ml_api.SCALER_MEAN[feature]
            ) / ml_api.SCALER_SCALE[feature]

        np.testing.assert_allclose(scaled, expected)
        np.testing.assert_allclose(X, np.array([[1.0, 30.0, 5.0, 8.0, 7.0]]))

    def test_resolve_model_path_uses_first_existing_candidate(self):
        with tempfile.NamedTemporaryFile(suffix=".joblib") as model_file:
            missing = str(Path(model_file.name).with_name("missing.joblib"))
            ml_api.CANDIDATES = [missing, model_file.name]

            resolved = ml_api.resolve_model_path()

        self.assertEqual(resolved, Path(model_file.name).resolve())

    def test_resolve_model_path_reports_missing_candidates(self):
        ml_api.CANDIDATES = ["missing-a.joblib", "missing-b.joblib"]

        with self.assertRaises(RuntimeError) as ctx:
            ml_api.resolve_model_path()

        self.assertIn("missing-a.joblib", str(ctx.exception))
        self.assertIn("missing-b.joblib", str(ctx.exception))

    def test_healthz_returns_ok(self):
        self.assertEqual(ml_api.healthz(), {"status": "ok"})

    def test_features_endpoint_returns_names_and_count(self):
        self.set_fake_artifact(FakeModel(at_risk_probability=0.42))

        response = ml_api.features()

        self.assertEqual(response["feature_names"], FEATURE_NAMES)
        self.assertEqual(response["count"], len(FEATURE_NAMES))

    def test_predict_uses_at_risk_probability_for_risk_percent(self):
        for probability, category in ((0.2, "low"), (0.6, "medium"), (0.9, "high")):
            with self.subTest(probability=probability, category=category):
                model = FakeModel(at_risk_probability=probability)
                self.set_fake_artifact(model)

                response = ml_api.predict(
                    ml_api.PredictRequest(features=self.valid_features())
                )

                self.assertAlmostEqual(response.RiskPercent, probability)
                self.assertEqual(response.Category, category)
                self.assertTrue(response.Message)

    def test_predict_uses_weighted_cascade_risk_when_second_stage_exists(self):
        model1 = FakeModel(at_risk_probability=0.8)
        model2 = FakeModel(at_risk_probability=0.25)
        self.set_fake_artifact(model1, model2=model2)

        response = ml_api.predict(ml_api.PredictRequest(features=self.valid_features()))

        self.assertAlmostEqual(response.RiskPercent, 0.7)
        self.assertEqual(response.Category, "medium")

    def test_prepare_X_for_model_scales_raw_request_features(self):
        features = self.valid_features()

        X = ml_api.prepare_X_for_model(features, FEATURE_NAMES)

        expected = ml_api.apply_training_scaling(
            ml_api.build_X(features, FEATURE_NAMES),
            FEATURE_NAMES,
        )
        np.testing.assert_allclose(X, expected)

    def test_predict_scales_raw_request_features_before_model_inference(self):
        model = FakeModel(at_risk_probability=0.6)
        self.set_fake_artifact(model)

        ml_api.predict(ml_api.PredictRequest(features=self.valid_features()))

        expected = ml_api.apply_training_scaling(
            ml_api.build_X(self.valid_features(), FEATURE_NAMES),
            FEATURE_NAMES,
        )
        np.testing.assert_allclose(model.seen_X, expected)

    def test_predict_rejects_missing_features(self):
        self.set_fake_artifact(FakeModel(at_risk_probability=0.6))

        with self.assertRaises(HTTPException) as ctx:
            ml_api.predict(ml_api.PredictRequest(features={"BMI": 30.0}))

        self.assertEqual(ctx.exception.status_code, 400)
        self.assertEqual(
            ctx.exception.detail,
            {
                "error": "Missing features",
                "missing": ["HighBP", "MentHlth", "PhysHlth", "Age"],
            },
        )


if __name__ == "__main__":
    unittest.main()
