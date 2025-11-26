def median(values):
    sorted_values = sorted(values)
    n = len(sorted_values)
    if n == 0:
        return 0
    if n % 2 == 0:
        mid1 = sorted_values[n // 2 - 1]
        mid2 = sorted_values[n // 2]
        return (mid1 + mid2) / 2
    else:
        return sorted_values[n // 2]


def std_deviation(values):
    if len(values) == 0:
        return 0
    mean = sum(values) / len(values)
    squared_diffs = []
    for x in values:
        squared_diffs.append((x - mean) ** 2)
    variance = sum(squared_diffs) / len(values)
    return variance ** 0.5
