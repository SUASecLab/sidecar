FROM httpd:2.4

# Install required packages
RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get install make m4 wget tar -y
RUN apt-get install python3 python3-pip python-is-python3 -y
RUN apt-get install apache2-dev -y

# Install wsgi
RUN mkdir wsgi && \
    cd wsgi && \
    wget https://github.com/GrahamDumpleton/mod_wsgi/archive/refs/tags/4.9.3.tar.gz && \
    tar xfz 4.9.3.tar.gz && \
    cd mod_wsgi* && \
    ./configure && \
    make && \
    make install

# Symlink wsgi module
RUN ln -s /usr/lib/apache2/modules/mod_wsgi.so /usr/local/apache2/modules/mod_wsgi.so

# Copy program files
RUN mkdir -p /usr/local/wsgi
COPY . /usr/local/wsgi

# Adjust httpd config
RUN cat /usr/local/wsgi/config_additions.txt >> /usr/local/apache2/conf/httpd.conf

# Install dependencies
RUN pip install --upgrade pip
RUN pip install --no-cache-dir -r /usr/local/wsgi/requirements.txt
RUN cd /usr/local/wsgi/sidecar_lib && pip install .

# Cleanup
RUN apt-get remove wget apache2-dev m4 make -y
RUN apt-get autoremove -y

# Fixup permissons
RUN chown www-data:www-data /usr/local/apache2 -R
RUN chown www-data:www-data /usr/local/wsgi -R

# Change to non-root user
USER www-data
