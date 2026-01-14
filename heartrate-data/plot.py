import pandas as pd
import plotly.graph_objects as go

df = pd.read_csv("heartrate.csv")

df["time"] = pd.to_datetime(df["measured_at_ms"], unit="ms")

df = df.sort_values("time")

# Filter to first hour only
start_time = df["time"].iloc[0]
end_time = start_time + pd.Timedelta(hours=1)
df = df[df["time"] <= end_time].copy()

df["heart_rate_smooth"] = df["heart_rate"].rolling(window=5, min_periods=1).mean()

# spike detection
window_seconds = 10
spike_threshold = 5
window_rows = window_seconds

df["heart_rate_window_diff"] = df["heart_rate_smooth"].rolling(
    window=window_rows, min_periods=1
).apply(lambda x: x[-1] - x[0], raw=True)

df["spike"] = df["heart_rate_window_diff"] >= spike_threshold
df["spike_value"] = df["heart_rate_smooth"].where(df["spike"])


fig = go.Figure()

# Raw data
fig.add_trace(go.Scatter(
    x=df["time"],
    y=df["heart_rate"],
    mode="markers",
    name="Raw",
    marker=dict(size=3, color="blue")
))

# Smoothed line
fig.add_trace(go.Scatter(
    x=df["time"],
    y=df["heart_rate_smooth"],
    mode="lines",
    name="Smoothed",
    line=dict(color="red", width=2)
))

# Rapid growth spike
fig.add_trace(go.Scatter(
    x=df["time"],
    y=df["spike_value"],
    mode="markers",
    name="Rapid increase",
    marker=dict(size=6, color="orange", symbol="triangle-up")
))

fig.update_layout(
    title="Heart Rate Over Time with Rapid Increases (10-second window)",
    xaxis_title="Time",
    yaxis_title="Heart Rate (bpm)",
    hovermode="x unified",
    xaxis=dict(
        rangeslider=dict(visible=True),
        type="date"
    )
)

fig.show()
