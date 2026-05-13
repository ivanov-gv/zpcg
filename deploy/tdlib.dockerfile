FROM golang:1.26-bookworm

ARG TDLIB_COMMIT=971684a

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        zlib1g-dev libssl-dev gperf php-cli cmake g++ libstdc++-dev && \
    git clone https://github.com/tdlib/td.git /tmp/td && \
    cd /tmp/td && git checkout ${TDLIB_COMMIT} && \
    mkdir build && cd build && \
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/usr/local .. && \
    cmake --build . --target install -j$(nproc) && \
    rm -rf /tmp/td && \
    apt-get purge -y cmake g++ php-cli gperf && \
    apt-get autoremove -y && rm -rf /var/lib/apt/lists/*

ENV CGO_CFLAGS="-I/usr/local/include" \
    CGO_LDFLAGS="-L/usr/local/lib -ltdjson" \
    LD_LIBRARY_PATH="/usr/local/lib"