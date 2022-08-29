FROM python:3-alpine

RUN apk -U upgrade
RUN apk add make m4

RUN adduser -D sidecar

WORKDIR /usr/src/app
COPY --chown=sidecar:sidecar . .

USER sidecar

ENV PATH=$PATH:/home/sidecar/.local/bin

RUN pip install --user --upgrade pip
RUN pip install --no-cache-dir --user -r requirements.txt
RUN make

CMD ["python", "-u", "main.py"]
