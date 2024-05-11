import json
import subprocess
from pyvis.network import Network
from IPython.core.display import display, HTML

class Node:
  def __init__(self, name, gateway, ip, mac, id):
    self.name = name
    self.gateway = gateway
    self.ip = ip
    self.mac = mac
    self.id = id

class AccessPoint:
  def __init__(self, name, driver, subnet, gateway, id, scope):  
    self.name = name
    self.driver = driver
    self.subnet = subnet
    self.gateway = gateway
    self.id = id
    self.scope = scope

class Graph:
  def __init__(self):
    self.access_point = []
    self.node = []

  def get_access_point(self):
    proc = subprocess.Popen(["docker", "network", "ls", "--format", "{{.ID}}/{{.Name}}/{{.Driver}}/{{.Scope}}"], stdout=subprocess.PIPE)
    (all_networks, err) = proc.communicate()
    all_networks = all_networks.decode().split("\n")
    all_networks.pop(len(all_networks) - 1)

    count = 0

    for network in all_networks:
      network = network.split("/")
      proc = subprocess.Popen(["docker", "network", "inspect", network[0]], stdout=subprocess.PIPE)
      (detail, err) = proc.communicate()
      detail = json.loads(detail.decode())[0]
      temp = AccessPoint(network[1], network[2], detail["IPAM"]["Config"][0]["Subnet"], detail["IPAM"]["Config"][0]["Gateway"], network[0], network[3])
      self.access_point.append(temp)
      count += 1
      if (count >= len(all_networks) - 2): break
      
  def get_node(self, docker_images):
    for di in docker_images:
      proc = subprocess.Popen(["docker", "inspect", di[7]], stdout=subprocess.PIPE)
      (detail, err) = proc.communicate()
      detail = json.loads(detail.decode())[0]
      temp = Node(di[2], detail["NetworkSettings"]["Gateway"],  detail["NetworkSettings"]["IPAddress"],  detail["NetworkSettings"]["MacAddress"], di[7])
      self.node.append(temp)

  def generate(self):
    net = Network(height="1500px",width="100%",bgcolor="#222222", font_color="white", directed=True)
    
    for node in self.node:
      net.add_node(node.name, label=node.name, shape="image", image="./static/Docker-Logo.png", title="ID: "+node.id+"\n\n"+"Gateway: "+node.gateway+"\n"+"IP: "+node.ip+"\n"+"MAC: "+node.mac)
    
    for access_point in self.access_point:
      net.add_node(access_point.name, label=access_point.name, shape="square", color="green", title="ID: "+access_point.id+"\n\n"+"Driver: "+access_point.driver+"\n"+"Scope: "+access_point.scope+"\n\n"+"Subnet: "+access_point.subnet+"\n"+"Gateway: "+access_point.gateway)
      for node in self.node:
        if node.gateway == access_point.gateway:
          net.add_edge(access_point.name, node.name)

    net.repulsion()
    net.save_graph("./templates/nodes.html")
    # display(HTML("graph.html"))

# a = Graph()
# a.get_access_point()
# print(a.access_point[0].name)

# docker_images = [[1, 'UP', 'beaver-yo-yuh', 'latest', '4722b5034e14', '576MB', '2024-04-15 20:42:05 +0700 +07', 'dd85f674cabf', 'Up 46 minutes', '2004,2005']]

# a = Graph()
# a.get_access_point()
# a.get_node(docker_images)
# a.generate()

