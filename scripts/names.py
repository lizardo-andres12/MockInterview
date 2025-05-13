from typing import Any, Generator
import sys


def get_name_from_file(filepath: str) -> Generator[str, Any, Any]:
    with open(filepath, 'r') as f:
        for line in f:
            yield line


def main():
    if len(sys.argv) < 2 :
        # print(usage)
        sys.exit(1)

    filepath = sys.argv[1]
    res = 'my_list=('

    try:
        for line in get_name_from_file(filepath):
            line = line[:len(line) - 1]
            res += ''.join(line.split(' ')) + ' '

    except FileNotFoundError:
        print(f'Error: no such file "{filepath}"')
        sys.exit(2)

    print(f'{res[:len(res) - 1]})')


if __name__ == '__main__':
    main()

