import multiprocessing
import os

bind = f"0.0.0.0:{os.getenv('PORT', 5001)}"
worker = multiprocessing.cpu_count() * 2 + 1
log_file = "-"
max_requests = 5000
max_requests_jitter = 50
access_log = "-"
