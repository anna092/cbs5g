import requests
import xml.etree.ElementTree as ET
from datetime import datetime, timezone, timedelta
import argparse

parser = argparse.ArgumentParser(description='The parameter for the CBS message')
parser.add_argument('-id', '--messageId', type=int, help='The message ID for serial number', required = True)
xml_data = """<?xml version="1.0" encoding="UTF-8"?> 
<alert xmlns="urn:oasis:names:tc:emergency:cap:1.1">
    <identifier>CWB-EQ112202</identifier>  
    <sender>cwb@scman.cwb.gov.tw</sender> 
    <sent>2023-07-27T00:08:00+08:00</sent> 
    <status>Actual</status> 
    <msgType>Alert</msgType>
    <source>CWB</source>
    <scope>Public</scope> 
    <info> 
        <language>zh-TW</language>
        <category>Met</category>
        <event>地震</event>
        <responseType>Shelter</responseType>
        <urgency>Immediate</urgency>
        <severity>Severe</severity>
        <certainty>Observed</certainty>
        <effective>2023-07-27T08:00:00+08:00</effective>
        <expires>2023-07-27T08:08:00+08:00</expires> 
        <senderName>中央氣象局</senderName> 
        <headline>地震報告</headline> 
        <description>花蓮縣秀林鄉發生規模 5.3 有感地震，最大震度花蓮縣太魯閣、宜蘭縣南山、南投縣合歡山、臺中市德基 4 級。</description>
        <contact>123456</contact>  
        <area> 
            <areaDesc>最大震度 3 級地區</areaDesc> 
            <geoCode>10002</geoCode> 
        </area> 
        </info> 
    </alert>
"""

args = parser.parse_args()
root = ET.fromstring(xml_data)
current_time_utc = datetime.now(timezone.utc)
taiwan_timezone = timezone(timedelta(hours=8))
taiwan_time = current_time_utc.astimezone(taiwan_timezone) 
formatted_time = taiwan_time.strftime('%Y-%m-%d %H:%M:%S.%f')[:-3] + ' ' + taiwan_time.tzname()
effective_element = root.find('.//{urn:oasis:names:tc:emergency:cap:1.1}effective')
if effective_element is not None:
            effective_element.text = formatted_time
sent_element = root.find('.//{urn:oasis:names:tc:emergency:cap:1.1}sent')
if sent_element is not None:
            sent_element.text = formatted_time
serialNumber_element = root.find('.//{urn:oasis:names:tc:emergency:cap:1.1}identifier')
if serialNumber_element is not None:
            serialNumber_element.text = serialNumber_element.text[:-3] + f"{args.messageId:03d}"
modified_xml_string = ET.tostring(root, encoding='utf-8', method='xml')

def send_xml_data(url, xml_data):
    headers = {
        "Content-Type": "application/xml ;charset=utf-8",
    }
    
    try:
        response = requests.post(url, data=xml_data, headers=headers)
        response.raise_for_status()
        return response.content
    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")
        return None

url = 'http://127.0.0.1:8080'

response_content = send_xml_data(url, modified_xml_string)
print("Data send at", formatted_time)
if response_content is None:
    print("Failed to send XML data.")
