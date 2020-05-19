#!/usr/bin/env python3

import os
import pandas as pd
import numpy as np
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
    tmp = dict()
    for zone in iter_files():
        with open(zone) as f:
            content = f.read().splitlines()
            llen = len(content)
            path = zone.split("/")[3:-1][::-1]
            while len(path) < 4:
                path += [None]
            print([llen]+path)
            df.loc[len(df)+1] = [llen]+path

def main():
    global err
    global df
    err = 0
    df = pd.DataFrame(columns=("len","tld","dom","sub", "ssub"))
    read_data()
    print(df)
    fig = px.sunburst(df,
                  title="some....",
                  #path=['dom','tld'],
                  path=['tld','dom'],
                  values='len',
                  #color='len',
                  hover_data=['len'],
                  color_continuous_scale='RdBu',
                  color_continuous_midpoint=np.average(df['len'],
                  weights=df['len']))
    fig.show()

if __name__ == "__main__":
    main()
