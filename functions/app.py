# app.py
import sys
import time
def main():
    # Leer argumento de línea de comandos
    if len(sys.argv) < 2:
        print("Error: No parameter provided")
        sys.exit(1)

    param = sys.argv[1]
    # Procesar la entrada
    result = param[::-1]
    
    # Simular una ejecución lenta con un retraso de 5 segundos
    time.sleep(20)
    print(result)  # Imprimir resultado al stdout

if __name__ == "__main__":
    main()
