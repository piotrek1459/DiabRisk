import unittest
from pathlib import Path
from unittest.mock import patch

import numpy as np
import pandas as pd

from src.Ml import main as ml_training


class MLTrainingTests(unittest.TestCase):
    def processed_frames(self, x_test_columns=("A", "B")):
        return [
            pd.DataFrame({"A": [1.5, 3.5], "B": [2.5, 4.5]}),
            pd.DataFrame({x_test_columns[0]: [5.5], x_test_columns[1]: [6.5]}),
            pd.DataFrame({"Diabetes_012": [0.0, 2.0]}),
            pd.DataFrame({"Diabetes_012": [1.0]}),
        ]

    def test_load_processed_splits_reads_processed_train_and_test_files(self):
        with patch.object(
            ml_training.pd,
            "read_csv",
            side_effect=self.processed_frames(),
        ) as read_csv:
            X_train, X_test, y_train, y_test, feature_names = (
                ml_training.load_processed_splits(Path("repo"))
            )

        self.assertEqual(feature_names, ["A", "B"])
        np.testing.assert_allclose(X_train, np.array([[1.5, 2.5], [3.5, 4.5]]))
        np.testing.assert_allclose(X_test, np.array([[5.5, 6.5]]))
        np.testing.assert_array_equal(y_train, np.array([0, 2]))
        np.testing.assert_array_equal(y_test, np.array([1]))
        self.assertEqual(
            [call.args[0].name for call in read_csv.call_args_list],
            [
                "X_train_processed.csv",
                "X_test_processed.csv",
                "y_train.csv",
                "y_test.csv",
            ],
        )

    def test_load_processed_splits_rejects_mismatched_feature_columns(self):
        with patch.object(
            ml_training.pd,
            "read_csv",
            side_effect=self.processed_frames(x_test_columns=("A", "C")),
        ):
            with self.assertRaises(ValueError) as ctx:
                ml_training.load_processed_splits(Path("repo"))

        self.assertEqual(
            str(ctx.exception),
            "Processed train and test feature columns do not match",
        )


if __name__ == "__main__":
    unittest.main()
