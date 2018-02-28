import random as rand

keys = []

for i in range(0,10000):
    key = rand.randint(0,10000)
    keys.append(key)
    print("put", 1, key, rand.randint(0, 10000))

print()
for i in range(0,10000):
    print("get", 1, keys[i])
