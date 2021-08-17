FROM debian:stable

LABEL maintainer "cherzog@web.de"

RUN apt-get update && \
    apt-get install -y golang git make && \
    git clone https://github.com/cherzog/csvdiff.git && \
    cd csvdiff && \
    make && \
    make install && \
    cp /csvdiff/out/csvdiff /csv-diff && \
    rm -R /csvdiff && \
    apt-get remove -y golang git make && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

CMD ["/csv-diff", "/b", "/d"]
#CMD ["/bin/bash"]
