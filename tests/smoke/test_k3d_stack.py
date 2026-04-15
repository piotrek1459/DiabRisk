import json
import os
import platform
import subprocess
import time
import unittest
import urllib.error
import urllib.parse
import urllib.request
from http.cookiejar import CookieJar


BASE_URL = os.environ.get("DIABRISK_BASE_URL", "http://localhost")
CLUSTER_NAME = os.environ.get("DIABRISK_K3D_CLUSTER", "diabrisk")
SKIP_CLUSTER_SETUP = os.environ.get("DIABRISK_SKIP_K3D_SETUP") == "1"
DELETE_CLUSTER_AFTER = os.environ.get("DIABRISK_DELETE_K3D_AFTER") == "1"


class HTTPResponseError(AssertionError):
    pass


def run_command(command, timeout=900):
    completed = subprocess.run(
        command,
        text=True,
        encoding="utf-8",
        errors="replace",
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        timeout=timeout,
    )
    if completed.returncode != 0:
        raise AssertionError(
            f"Command failed with exit code {completed.returncode}: {' '.join(command)}\n"
            f"{completed.stdout}"
        )
    return completed.stdout


def tail(text, max_lines=12):
    lines = text.strip().splitlines()
    return "\n".join(lines[-max_lines:])


def require_command(command, hint):
    try:
        run_command(command, timeout=30)
    except FileNotFoundError as exc:
        raise AssertionError(f"{hint}\nCommand not found: {command[0]}") from exc
    except subprocess.TimeoutExpired as exc:
        raise AssertionError(f"{hint}\nCommand timed out: {' '.join(command)}") from exc
    except AssertionError as exc:
        raise AssertionError(f"{hint}\n{tail(str(exc))}") from exc


def check_cluster_prerequisites():
    require_command(["docker", "info"], "Docker is required for k3d smoke tests. Start Docker Desktop and try again.")
    require_command(["k3d", "version"], "k3d is required for full-stack smoke tests.")
    require_command(["kubectl", "version", "--client"], "kubectl is required for full-stack smoke tests.")


def install_cluster_command():
    if platform.system() == "Windows":
        return [
            "powershell.exe",
            "-ExecutionPolicy",
            "Bypass",
            "-File",
            "scripts\\install-local-k3d.ps1",
            "-ClusterName",
            CLUSTER_NAME,
        ]

    return ["bash", "scripts/install-local-k3d.sh", CLUSTER_NAME]


def delete_cluster_command():
    return ["k3d", "cluster", "delete", CLUSTER_NAME]


def wait_for_http(url, timeout=120):
    deadline = time.time() + timeout
    last_error = None
    while time.time() < deadline:
        try:
            with urllib.request.urlopen(url, timeout=5) as response:
                if response.status < 500:
                    return
        except OSError as exc:
            last_error = exc
        time.sleep(2)

    raise AssertionError(f"{url} did not become reachable: {last_error}")


def request_json(opener, method, path, payload=None, expected_status=200):
    data = None
    headers = {"Accept": "application/json"}
    if payload is not None:
        data = json.dumps(payload).encode("utf-8")
        headers["Content-Type"] = "application/json"

    request = urllib.request.Request(
        urllib.parse.urljoin(BASE_URL, path),
        data=data,
        headers=headers,
        method=method,
    )

    try:
        with opener.open(request, timeout=20) as response:
            body = response.read().decode("utf-8")
            status = response.status
    except urllib.error.HTTPError as exc:
        body = exc.read().decode("utf-8")
        status = exc.code

    if status != expected_status:
        raise HTTPResponseError(
            f"{method} {path} returned {status}, expected {expected_status}. Body: {body}"
        )

    return json.loads(body) if body else {}


def feature_payload(feature_names):
    defaults = {feature: 0 for feature in feature_names}
    defaults.update(
        {
            "HighBP": 1,
            "HighChol": 1,
            "CholCheck": 1,
            "BMI": 30.0,
            "Smoker": 0,
            "Stroke": 0,
            "HeartDiseaseorAttack": 0,
            "PhysActivity": 1,
            "Fruits": 1,
            "Veggies": 1,
            "HvyAlcoholConsump": 0,
            "AnyHealthcare": 1,
            "NoDocbcCost": 0,
            "GenHlth": 3,
            "MentHlth": 5,
            "PhysHlth": 3,
            "DiffWalk": 0,
            "Sex": 1,
            "Age": 9,
            "Education": 5,
            "Income": 6,
        }
    )
    return {feature: defaults[feature] for feature in feature_names}


class K3DStackSmokeTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        if not SKIP_CLUSTER_SETUP:
            check_cluster_prerequisites()
            run_command(install_cluster_command())
        wait_for_http(BASE_URL)

    @classmethod
    def tearDownClass(cls):
        if DELETE_CLUSTER_AFTER:
            run_command(delete_cluster_command(), timeout=180)

    def setUp(self):
        cookie_jar = CookieJar()
        self.opener = urllib.request.build_opener(
            urllib.request.HTTPCookieProcessor(cookie_jar)
        )

    def test_frontend_is_served(self):
        with urllib.request.urlopen(BASE_URL, timeout=20) as response:
            body = response.read().decode("utf-8", errors="replace").lower()

        self.assertEqual(response.status, 200)
        self.assertIn("<html", body)

    def test_auth_gateway_ml_happy_path(self):
        unique = int(time.time() * 1000)
        email = f"smoke-{unique}@example.test"
        password = "smoke-password-123"

        register_response = request_json(
            self.opener,
            "POST",
            "/auth/register",
            {
                "email": email,
                "password": password,
                "full_name": "Smoke Test",
            },
            expected_status=201,
        )
        self.assertEqual(register_response["user"]["email"], email)

        session_response = request_json(
            self.opener,
            "GET",
            "/auth/session",
            expected_status=200,
        )
        self.assertEqual(session_response["email"], email)

        features_response = request_json(
            self.opener,
            "GET",
            "/api/features",
            expected_status=200,
        )
        feature_names = features_response["feature_names"]
        self.assertGreater(features_response["count"], 0)
        self.assertEqual(features_response["count"], len(feature_names))

        risk_response = request_json(
            self.opener,
            "POST",
            "/api/risk",
            {"features": feature_payload(feature_names)},
            expected_status=200,
        )

        self.assertIn("RiskPercent", risk_response)
        self.assertIn("Category", risk_response)
        self.assertIn("Message", risk_response)
        self.assertGreaterEqual(risk_response["RiskPercent"], 0)
        self.assertLessEqual(risk_response["RiskPercent"], 1)
        self.assertIn(risk_response["Category"], ["low", "medium", "high"])


if __name__ == "__main__":
    unittest.main()
