import multiprocessing

bind = "0.0.0.0:8080"
worker_class = 'gevent'
workers = multiprocessing.cpu_count() * 2 + 1
log_file = "-"
max_requests = 5000
max_requests_jitter = 50
access_log = "-"
