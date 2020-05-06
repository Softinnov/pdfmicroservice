FROM ubuntu:14.04

RUN sed 's/main$/main universe/' -i /etc/apt/sources.list
RUN apt-get update

# Install wkhtmltopdf
RUN apt-get install -y build-essential xorg libssl-dev libxrender-dev poppler-utils fontconfig xfonts-75dpi
ADD wkhtmltox-0.12.2.1_linux-trusty-amd64.deb /tmp/
RUN dpkg -i /tmp/wkhtmltox-0.12.2.1_linux-trusty-amd64.deb
RUN rm /tmp/wkhtmltox-0.12.2.1_linux-trusty-amd64.deb

EXPOSE 8000

ADD pdfmicroservice /share/pdfmicro
WORKDIR /share

VOLUME /share/pdf

CMD "/share/pdfmicro"
