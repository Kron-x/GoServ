import subprocess

def start_go_server():
    # Start Go-service
    subprocess.Popen(["go", "run", "main.go"])

if __name__ == "__main__":
    start_go_server()
    print("Go server is running...")
