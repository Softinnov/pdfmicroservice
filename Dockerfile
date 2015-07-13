FROM ubuntu:14.04

RUN sed 's/main$/main universe/' -i /etc/apt/sources.list
RUN apt-get update

# Download and install wkhtmltopdf
RUN apt-get install -y build-essential xorg libssl-dev libxrender-dev wget gdebi poppler-utils
RUN wget http://download.gna.org/wkhtmltopdf/0.12/0.12.2.1/wkhtmltox-0.12.2.1_linux-trusty-amd64.deb
RUN gdebi --n wkhtmltox-0.12.2.1_linux-trusty-amd64.deb

EXPOSE 8000

ADD pdfmicroservice /share/pdfmicro
WORKDIR /share

VOLUME /share/pdf

CMD "/share/pdfmicro"
