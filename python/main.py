import subprocess
import os

def start_go_server():
    # Получаем полный путь к директории с этим скриптом
    script_dir = os.path.dirname(os.path.abspath(__file__))
    # Переход в папку go
    os.chdir(script_dir)
    os.chdir("../go")
    # Запуск Go-сервера
    subprocess.Popen(["go", "run", "main.go"])

if __name__ == "__main__":
    start_go_server()
    print("Go server is running...")