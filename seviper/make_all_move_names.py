import boa
import base_data

def main():
    boa.write_txt(base_data.ALL_MOVE_NAMES_PATH,
                  "\n".join(base_data.ALL_MOVE_NAMES))

if __name__ == "__main__":
    main()
