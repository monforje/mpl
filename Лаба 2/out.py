from __future__ import annotations
from dataclasses import dataclass

@dataclass
class Counter:
    Value: int
    def Inc(self, step):
        self.Value = self.Value + step
        return self


def divmod(a, b):
    q = a // b
    r = a % b
    return q, r

def main():
    c = Counter(Value=0)
    i = 0
    while i < 5:
        c = c.Inc(2)
        i = i + 1
    q, r = divmod(c.Value, 3)
    if r == 0:
        print(q)
    else:
        print(q, r)

if __name__ == "__main__":
    main()
