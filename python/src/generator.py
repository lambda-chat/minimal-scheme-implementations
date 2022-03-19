from collections.abc import Callable, Generator
import sys


def generate_int(next: Callable[[int], int], x: int) -> Generator[int, None, None]:
    while True:
        yield x
        x = next(x)


def generator_main_1():
    gen = generate_int(lambda x: x + 1, 0)
    for _ in range(10):
        print(next(gen))


def client_process() -> Generator[int, int, str]:
    msg = yield 1
    print(f"[client] recieved {msg}", file=sys.stderr)
    msg = yield 2
    print(f"[client] recieved {msg}", file=sys.stderr)
    msg = yield 3
    print(f"[client] recieved {msg}", file=sys.stderr)
    return "client process done."


def server_process() -> Generator[int, int, str]:
    msg = yield 0  # discarded
    print(f"[server] recieved {msg}", file=sys.stderr)
    msg = yield 11
    print(f"[server] recieved {msg}", file=sys.stderr)
    msg = yield 22
    print(f"[server] recieved {msg}", file=sys.stderr)
    _ = yield 33
    return "server process done."


def generator_main_2() -> None:
    client = client_process()
    server = server_process()
    client_msg = next(client)
    server_msg = next(server)
    client_alive = True
    server_alive = True
    while client_alive or server_alive:
        try:
            if server_alive:
                server_msg = server.send(client_msg)
        except StopIteration as ex:
            msg: str = ex.value
            print(msg, file=sys.stderr)
            server_alive = False

        try:
            if client_alive:
                client_msg = client.send(server_msg)
        except StopIteration as ex:
            msg: str = ex.value
            print(msg, file=sys.stderr)
            client_alive = False


if __name__ == "__main__":
    generator_main_2()
