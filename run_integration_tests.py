"""
Simple test runner for integration tests.
Run this to execute all integration tests.
"""
import sys
import subprocess
from pathlib import Path

def run_tests():
    """Run all integration tests."""
    project_root = Path(__file__).parent.parent
    
    print("Running Integration Tests...")
    print("=" * 50)

    # Run API integration tests
    print("\n2. Running API Integration Tests...")
    result = subprocess.run([
        sys.executable, "-m", "pytest", 
        "tests/test_measurement_api.py", 
        "-v", "--tb=short"
    ], cwd=project_root)
    
    # Summary
    print("\n" + "=" * 50)
    if result.returncode == 0:
        print("✅ All integration tests PASSED!")
        return True
    else:
        print("❌ Some integration tests FAILED!")
        return False

if __name__ == "__main__":
    success = run_tests()
    sys.exit(0 if success else 1)
