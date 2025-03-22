import subprocess
import os

def start_go_server():
    # Переход в папку go
    os.chdir("../go")
    # Запуск Go-сервера
    subprocess.Popen(["go", "run", "main.go"])

if __name__ == "__main__":
    start_go_server()
    print("Go server is running...")
