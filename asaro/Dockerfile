###########
# BUILDER #
###########
FROM python:3.11.1-slim-buster as builder

WORKDIR /opt/app

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

RUN pip install --upgrade pip
COPY . .

RUN pip wheel --no-cache-dir --no-deps --wheel-dir /opt/wheels -r requirements.txt


#########
# FINAL #
#########
FROM python:3.11.1-slim-buster as final

RUN mkdir -p /opt/app

COPY --from=builder /opt/wheels /wheels
COPY --from=builder /opt/app .

RUN pip install --upgrade pip
RUN pip install --no-cache /wheels/*

CMD ["gunicorn", "app:application"]
