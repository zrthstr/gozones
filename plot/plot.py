#!/usr/bin/env python3

import os
import pandas as pd
import plotly.express as px

DATADIR = "../data/zones"

"""
data = dict(
    character=["Eve", "Cain", "Seth", "Enos", "Noam", "Abel", "Awan", "Enoch", "Azura"],
    parent=["", "Eve", "Eve", "Seth", "Seth", "Eve", "Eve", "Awan", "Eve" ],
    value=[10, 14, 12, 10, 2, 6, 6, 4, 4])

fig =px.sunburst(
    data,
    names='character',
    parents='parent',
    values='value',
)
"""

def iter_files():
    rootdir = DATADIR
    #for subdir, dirs, files in os.walk(rootdir):
    for subdir, _, files in os.walk(rootdir):
        for file in files:
            yield os.path.join(subdir, file)



def read_data():
    global err, df
    for zone in iter_files():
        print("zzzz",zone)
        with open(zone) as f:
            #data = np.array([['','Col1','Col2'], ['Row1',1,2], ['Row2',3,4]])
            content = f.read().splitlines()
            #print(content)
            clean_lines = []
            for line in content:
                try:
                    path = zone.split("/")[3:-1][::-1]
                    print("PATH:",path)
                    this = line.split("\t")
                    this.extend([zone])
                    clean_lines.append(this)
                    #clean_lines.extend(line.split("\t"))
                    TLD, TTL, IN, TYPE, RR = line.split("\t")
                    print("-"*40)
                    print("line:", line)
                    print("TLD", TLD)
                    print("TTL", TTL)
                    print("IN", IN)
                    print("TYPE", TYPE)
                    print("RR", RR)
                    print("-"*40)
                except:
                    err += 1
                    continue
            df = pd.concat([df, pd.DataFrame(clean_lines)])

def main():
    global err
    global df
    err = 0
    # columns=list('AB')
    df = pd.DataFrame()
    read_data()
    print("err: ", err)
    print(df)
    #fig.show()

if __name__ == "__main__":
    main()
