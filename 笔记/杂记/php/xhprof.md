#xhprof

---

生产环境不要使用.

		# Install XHProf RUN wget https://github.com/RustJason/xhprof/archive/php7.tar.gz && \
    tar zxvf php7.tar.gz && \
    rm -f php7.tar.gz && \
    cd xhprof-php7/extension/ && \
    phpize && \
    ./configure --with-php-config=/usr/bin/php-config7.0 && \
    make && \
    make install && \
    rm -Rf ../xhprof-php7