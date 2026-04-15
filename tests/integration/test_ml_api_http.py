import asyncio
import json
import tempfile
import unittest
from pathlib import Path

import numpy as np
from joblib import dump

from src.FastAPI import ml_api


FEATURE_NAMES = ["HighBP", "BMI", "MentHlth", "PhysHlth", "Age"]


class IntegrationModel:
    def predict_proba(self, X):
        return np.array([[0.37, 0.63]], dtype=float)


def asgi_request(app, method, path, payload):
    body = json.dumps(payload).encode("utf-8")
    messages = []
    request_sent = False

    scope = {
        "type": "http",
        "asgi": {"version": "3.0"},
        "http_version": "1.1",
        "method": method,
        "scheme": "http",
        "path": path,
        "raw_path": path.encode("ascii"),
        "query_string": b"",
        "headers": [
            (b"host", b"testserver"),
            (b"content-type", b"application/json"),
            (b"content-length", str(len(body)).encode("ascii")),
        ],
        "client": ("127.0.0.1", 123),
        "server": ("testserver", 80),
    }

    async def receive():
        nonlocal request_sent
        if request_sent:
            return {"type": "http.disconnect"}
        request_sent = True
        return {"type": "http.request", "body": body, "more_body": False}

    async def send(message):
        messages.append(message)

    asyncio.run(app(scope, receive, send))

    status = next(
        message["status"]
        for message in messages
        if message["type"] == "http.response.start"
    )
    response_body = b"".join(
        message.get("body", b"")
        for message in messages
        if message["type"] == "http.response.body"
    )

    return status, json.loads(response_body)


class MLApiHTTPIntegrationTests(unittest.TestCase):
    def setUp(self):
        self.original_artifact = ml_api._artifact
        self.original_model_path = ml_api.MODEL_PATH
        self.original_candidates = ml_api.CANDIDATES

        self.temp_dir = tempfile.TemporaryDirectory()
        model_path = Path(self.temp_dir.name) / "model.joblib"
        dump(
            {
                "model1": IntegrationModel(),
                "feature_names": FEATURE_NAMES,
                "model_features": "processed",
                "api_input_features": "raw",
            },
            model_path,
        )

        ml_api._artifact = None
        ml_api.MODEL_PATH = None
        ml_api.CANDIDATES = [str(model_path)]

    def tearDown(self):
        ml_api._artifact = self.original_artifact
        ml_api.MODEL_PATH = self.original_model_path
        ml_api.CANDIDATES = self.original_candidates
        self.temp_dir.cleanup()

    def test_predict_loads_model_artifact_and_returns_http_contract(self):
        status_code, response_body = asgi_request(
            ml_api.app,
            "POST",
            "/predict",
            {
                "features": {
                    "HighBP": 1,
                    "BMI": 30.0,
                    "MentHlth": 5,
                    "PhysHlth": 8,
                    "Age": 7,
                }
            },
        )

        self.assertEqual(status_code, 200)
        self.assertEqual(
            response_body,
            {
                "RiskPercent": 0.63,
                "Category": "medium",
                "Message": "Moderate risk detected. Consider scheduling a medical checkup.",
            },
        )


if __name__ == "__main__":
    unittest.main()
