import matplotlib.pyplot as plt
import numpy as np

# Data based on your clean runs
scenarios = ['Baseline (0% Drop)', 'Degraded (25% Drop)']
p99_latencies = [4.17, 4.16]
p95_latencies = [2.88, 2.90]

x = np.arange(len(scenarios))
width = 0.35

fig, ax = plt.subplots(figsize=(8, 6))

# Plot bars
rects1 = ax.bar(x - width/2, p99_latencies, width, label='P99 Latency (ms)', color='#2c3e50')
rects2 = ax.bar(x + width/2, p95_latencies, width, label='P95 Latency (ms)', color='#3498db')

# Styling
ax.set_ylabel('Latency (ms)')
ax.set_title('Controller Resilience: Impact of Observability Faults')
ax.set_xticks(x)
ax.set_xticklabels(scenarios)
ax.set_ylim(0, 5) # Setting limit to highlight the stability
ax.legend()

# Add values on top of bars
def autolabel(rects):
    for rect in rects:
        height = rect.get_height()
        ax.annotate(f'{height:.2f}', xy=(rect.get_x() + rect.get_width() / 2, height),
                    xytext=(0, 3), textcoords="offset points", ha='center', va='bottom')

autolabel(rects1)
autolabel(rects2)

plt.tight_layout()
plt.savefig('thesis_resilience_graph.pdf') # Save as PDF for LaTeX
print("Graph saved as thesis_resilience_graph.pdf")
