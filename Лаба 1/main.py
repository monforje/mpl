import os
import random
import pandas as pd
from concurrent.futures import ThreadPoolExecutor
from statistics import median, std_deviation


def generate_csv_files():
    os.makedirs('1_step', exist_ok=True)

    for i in range(1, 6):
        data = {'Категория': [random.choice(['A', 'B', 'C', 'D']) for _ in range(
            100)], 'Значение': [random.uniform(0, 100) for _ in range(100)]}
        df = pd.DataFrame(data)
        df.to_csv(f'1_step/data_{i}.csv', index=False)


def process_file(file_num):
    input_file = f'1_step/data_{file_num}.csv'
    output_file = f'2_step/data_{file_num}.1.csv'

    df = pd.read_csv(input_file)

    results = []
    for c in ['A', 'B', 'C', 'D']:
        c_data = df[df['Категория'] == c]['Значение'].tolist()
        if c_data:
            med = median(c_data)
            std = std_deviation(c_data)
            results.append({'Категория': c, 'Медиана': med, 'Отклонение': std})
    result_df = pd.DataFrame(results)
    result_df.to_csv(output_file, index=False)


def parallel_process_files():
    os.makedirs('2_step', exist_ok=True)

    with ThreadPoolExecutor() as executor:
        executor.map(process_file, range(1, 6))


def create_final_file():
    os.makedirs('3_step', exist_ok=True)

    all_data = {'A': {'medians': [], 'stds': []},
                'B': {'medians': [], 'stds': []},
                'C': {'medians': [], 'stds': []},
                'D': {'medians': [], 'stds': []}}

    for i in range(1, 6):
        df = pd.read_csv(f'2_step/data_{i}.1.csv')
        for _, row in df.iterrows():
            category = row['Категория']
            all_data[category]['medians'].append(row['Медиана'])
            all_data[category]['stds'].append(row['Отклонение'])

    final_results = []
    for category in ['A', 'B', 'C', 'D']:
        med_of_meds = median(all_data[category]['medians'])
        std_of_stds = std_deviation(all_data[category]['stds'])
        final_results.append({
            'Категория': category,
            'Медиана медиан': med_of_meds,
            'Отклонение отклонений': std_of_stds
        })

    final_df = pd.DataFrame(final_results)
    final_df.to_csv('3_step/data_final.csv', index=False)


if __name__ == '__main__':
    generate_csv_files()
    parallel_process_files()
    create_final_file()
