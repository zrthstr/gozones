#!/usr/bin/env python3

import plotly.express as px
import numpy as np
df = px.data.gapminder().query("year == 2007")

print(df)

fig = px.sunburst(df,
                  #title="some....",
                  path=['continent', 'country'],
                  values='pop',
                  color='lifeExp',
                  hover_data=['iso_alpha'],
                  color_continuous_scale='RdBu',
                  color_continuous_midpoint=np.average(df['lifeExp'],
                  weights=df['pop']))
fig.show()
