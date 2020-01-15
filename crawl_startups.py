# -*- coding: UTF-8 -*-
##########################################################################################
# 3Lines VC
# Startup List
# Retrieve FOUNDERS, INVESTORS, EMAIL, WEBSITE, LOCATION
# zaubacorp.com, hoovers.com, corporatedir.com
##########################################################################################

import urllib2
import logging
import os
from selenium import webdriver
from selenium.webdriver.common.desired_capabilities import DesiredCapabilities
import time
import requests
from bs4 import BeautifulSoup
import ast
import re


def format_data():
    with open('company_details', 'r') as file:
        lines = file.readlines()
        lines = lines[:100]
        for line in lines:
            line = ast.literal_eval(line)
            directors = line['Director Details']
            company = line[line.keys()[1]]
            d_id, d_name, desig, appointment_date = directors[0], directors[1], directors[2], directors[3]
            print list(filter(re.compile("Email ID*").match, company)), d_id, d_name, desig, appointment_date


def get_company_url(soup, startup, headers, company_details_file):
    """
    :param soup: html
    :type startup: startup
    """
    for table in soup.find_all('table'):
        for link in table.find_all('a'):
            if startup.strip().upper() in str(link['href']) or startup.strip().upper().replace(" ", "-") in str(
                    link['href']):
                html = requests.get(link['href'], headers=headers)
                html_text = BeautifulSoup(html.text, 'html.parser')
                get_company_details(html_text, link['href'], company_details_file)
                break


def get_company_details(html_text, link, company_details_file):
    s_dict = {}
    a = []
    b = []
    for div in html_text.find_all(class_='accordion-toggle main-row'):
        [b.append(str(value.text.encode('utf-8')).strip()) for value in div.find_all('p')]

    for div in html_text.find_all(class_='col-lg-6 col-md-6 col-sm-12 col-xs-12'):
        [a.append(str(value.text.encode('utf-8')).strip()) for value in div.find_all('p')]

    s_dict[link] = a
    s_dict['Director Details'] = b
    company_details_file.write(str(s_dict) + '\n')


def main():
    headers = {
        'User-Agent': 'Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Trident/5.0)'
    }

    company_details_file = open('company_details', 'w')
    startup_file = open("startup_list_formatted", "r")

    for startup in startup_file.readlines():
        url = 'https://www.zaubacorp.com/companysearchresults/{0}/'.format(startup.upper().replace(" ", "-"))
        time.sleep(2)
        source = requests.get(url, headers=headers)
        soup = BeautifulSoup(source.text, 'html.parser')

        if str(soup.find(class_='breadcrumb').encode('utf-8')) in "COMPANY NOT FOUND" and startup in "Pvt. Ltd":
            startup = startup.replace("Pvt. Ltd", "PRIVATE LIMITED")
            source = requests.get('https://www.zaubacorp.com/companysearchresults/{0}/'
                                  .format(startup.upper().replace(" ", "-")), headers=headers)
            soup = BeautifulSoup(source.text, 'html.parser')
            if str(soup.find(class_='breadcrumb').encode('utf-8')) in "COMPANY NOT FOUND":
                continue
            else:
                get_company_url(soup, startup, headers, company_details_file)
        else:
            get_company_url(soup, startup, headers, company_details_file)
    company_details_file.close()


if __name__ == "__main__":
    # main()
    format_data()
