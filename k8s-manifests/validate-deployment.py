#!/usr/bin/env python3
"""
E-commerce Kubernetes Validation Script
Tests the complete deployment on Kubernetes cluster
"""

import requests
import time
import subprocess
import json
from datetime import datetime

class KubernetesValidator:
    def __init__(self):
        self.minikube_ip = None
        self.services = {
            'catalog': {'port': 30001, 'health_path': '/health'},
            'inventory': {'port': 30002, 'health_path': '/health'},
            'customer': {'port': 30003, 'health_path': '/health'},
            'payment': {'port': 30004, 'health_path': '/health'},
            'order': {'port': 30005, 'health_path': '/docs'},
            'shipment': {'port': 30006, 'health_path': '/docs'},
            'notification': {'port': 30007, 'health_path': '/actuator/health'}
        }
        
    def print_step(self, step, title, emoji="🔄"):
        print(f"\n{emoji} {step}: {title}")
        print("-" * 60)
        
    def print_success(self, message):
        print(f"✅ {message}")
        
    def print_error(self, message):
        print(f"❌ {message}")
        
    def print_info(self, message):
        print(f"ℹ️  {message}")

    def get_minikube_ip(self):
        """Get Minikube IP address"""
        try:
            result = subprocess.run(['minikube', 'ip'], capture_output=True, text=True)
            if result.returncode == 0:
                self.minikube_ip = result.stdout.strip()
                self.print_success(f"Minikube IP: {self.minikube_ip}")
                return True
            else:
                self.print_error("Failed to get Minikube IP")
                return False
        except Exception as e:
            self.print_error(f"Error getting Minikube IP: {str(e)}")
            return False

    def check_kubernetes_pods(self):
        """Check if all pods are running"""
        self.print_step("1", "Checking Kubernetes Pods Status", "⚓")
        
        try:
            # Check ecommerce namespace pods
            result = subprocess.run(
                ['kubectl', 'get', 'pods', '-n', 'ecommerce', '-o', 'json'],
                capture_output=True, text=True
            )
            
            if result.returncode != 0:
                self.print_error("Failed to get pod status")
                return False
                
            pods_data = json.loads(result.stdout)
            running_pods = 0
            total_pods = len(pods_data['items'])
            
            for pod in pods_data['items']:
                pod_name = pod['metadata']['name']
                pod_status = pod['status']['phase']
                
                if pod_status == 'Running':
                    running_pods += 1
                    self.print_success(f"Pod {pod_name}: {pod_status}")
                else:
                    self.print_error(f"Pod {pod_name}: {pod_status}")
            
            self.print_info(f"Pods Running: {running_pods}/{total_pods}")
            return running_pods == total_pods
            
        except Exception as e:
            self.print_error(f"Error checking pods: {str(e)}")
            return False

    def check_kubernetes_services(self):
        """Check if all services are accessible"""
        self.print_step("2", "Checking Kubernetes Services", "🌐")
        
        try:
            result = subprocess.run(
                ['kubectl', 'get', 'services', '-n', 'ecommerce', '-o', 'json'],
                capture_output=True, text=True
            )
            
            if result.returncode != 0:
                self.print_error("Failed to get service status")
                return False
                
            services_data = json.loads(result.stdout)
            
            for service in services_data['items']:
                service_name = service['metadata']['name']
                service_type = service['spec']['type']
                
                if 'ports' in service['spec']:
                    ports = [str(port['port']) for port in service['spec']['ports']]
                    ports_str = ','.join(ports)
                    self.print_success(f"Service {service_name}: {service_type} (ports: {ports_str})")
                else:
                    self.print_success(f"Service {service_name}: {service_type}")
            
            return True
            
        except Exception as e:
            self.print_error(f"Error checking services: {str(e)}")
            return False

    def test_service_health(self):
        """Test health endpoints of all services"""
        self.print_step("3", "Testing Service Health Endpoints", "🏥")
        
        if not self.minikube_ip:
            self.print_error("Minikube IP not available")
            return False
            
        healthy_services = 0
        
        for service_name, config in self.services.items():
            try:
                url = f"http://{self.minikube_ip}:{config['port']}{config['health_path']}"
                response = requests.get(url, timeout=10)
                
                if response.status_code == 200:
                    self.print_success(f"{service_name.title()} Service: Healthy")
                    healthy_services += 1
                else:
                    self.print_error(f"{service_name.title()} Service: HTTP {response.status_code}")
                    
            except requests.exceptions.RequestException as e:
                self.print_error(f"{service_name.title()} Service: {str(e)}")
        
        total_services = len(self.services)
        self.print_info(f"Healthy Services: {healthy_services}/{total_services}")
        return healthy_services == total_services

    def test_api_functionality(self):
        """Test basic API functionality"""
        self.print_step("4", "Testing API Functionality", "🔧")
        
        if not self.minikube_ip:
            return False
            
        success_count = 0
        
        # Test Catalog Service - Create Product
        try:
            catalog_url = f"http://{self.minikube_ip}:30001/v1/products"
            product_data = {
                "sku": "K8S-TEST-001",
                "name": "Kubernetes Test Product",
                "category": "Testing",
                "price": 99.99,
                "is_active": True,
                "description": "Product created during K8s validation"
            }
            
            response = requests.post(catalog_url, json=product_data, timeout=10)
            if response.status_code in [200, 201]:
                self.print_success("Catalog API: Product creation successful")
                success_count += 1
            else:
                self.print_error(f"Catalog API: Failed ({response.status_code})")
                
        except Exception as e:
            self.print_error(f"Catalog API: {str(e)}")

        # Test Inventory Service - Add Stock
        try:
            inventory_url = f"http://{self.minikube_ip}:30002/v1/inventory"
            inventory_data = {
                "product_id": 1,
                "warehouse": "K8S-WH",
                "onhand": 100,
                "reserved": 0
            }
            
            response = requests.post(inventory_url, json=inventory_data, timeout=10)
            if response.status_code in [200, 201]:
                self.print_success("Inventory API: Stock addition successful")
                success_count += 1
            else:
                self.print_error(f"Inventory API: Failed ({response.status_code})")
                
        except Exception as e:
            self.print_error(f"Inventory API: {str(e)}")

        # Test Payment Service - Health Check (already tested above, but verify again)
        try:
            payment_url = f"http://{self.minikube_ip}:30004/health"
            response = requests.get(payment_url, timeout=10)
            if response.status_code == 200:
                self.print_success("Payment API: Health check successful")
                success_count += 1
            else:
                self.print_error(f"Payment API: Health check failed ({response.status_code})")
                
        except Exception as e:
            self.print_error(f"Payment API: {str(e)}")
        
        self.print_info(f"API Tests Passed: {success_count}/3")
        return success_count >= 2

    def check_monitoring_stack(self):
        """Check if monitoring stack is deployed"""
        self.print_step("5", "Checking Monitoring Stack", "📊")
        
        try:
            # Check monitoring namespace
            result = subprocess.run(
                ['kubectl', 'get', 'pods', '-n', 'monitoring', '-o', 'json'],
                capture_output=True, text=True
            )
            
            if result.returncode != 0:
                self.print_info("Monitoring stack not deployed (optional)")
                return True
                
            pods_data = json.loads(result.stdout)
            
            if len(pods_data['items']) == 0:
                self.print_info("Monitoring stack not deployed (optional)")
                return True
                
            monitoring_services = ['prometheus', 'grafana']
            found_services = []
            
            for pod in pods_data['items']:
                pod_name = pod['metadata']['name']
                for service in monitoring_services:
                    if service in pod_name.lower():
                        pod_status = pod['status']['phase']
                        if pod_status == 'Running':
                            self.print_success(f"Monitoring: {service.title()} is running")
                            found_services.append(service)
                        else:
                            self.print_error(f"Monitoring: {service.title()} is {pod_status}")
            
            if found_services:
                self.print_info(f"Monitoring URLs:")
                self.print_info(f"  • Prometheus: http://{self.minikube_ip}:30090")
                self.print_info(f"  • Grafana: http://{self.minikube_ip}:30030 (admin/admin123)")
            
            return True
            
        except Exception as e:
            self.print_info(f"Monitoring check failed: {str(e)} (optional)")
            return True

    def generate_summary_report(self):
        """Generate final validation report"""
        self.print_step("6", "Deployment Validation Summary", "📋")
        
        print(f"""
🎉 KUBERNETES DEPLOYMENT VALIDATION COMPLETE! 🎉

📊 Cluster Information:
   • Minikube IP: {self.minikube_ip}
   • Namespace: ecommerce
   • Services: {len(self.services)} microservices deployed

🔗 Service Access URLs:
""")
        
        for service_name, config in self.services.items():
            port = config['port']
            print(f"   • {service_name.title()} Service: http://{self.minikube_ip}:{port}")
        
        print(f"""
🛠️  Management Commands:
   • View pods: kubectl get pods -n ecommerce
   • View services: kubectl get services -n ecommerce  
   • View logs: kubectl logs deployment/<service-name> -n ecommerce
   • Scale service: kubectl scale deployment <service-name> --replicas=3 -n ecommerce
   • Port forward: kubectl port-forward service/<service-name> <local-port>:<service-port> -n ecommerce

🔧 Quick Tests:
   • Health check: curl http://{self.minikube_ip}:30001/health
   • Get products: curl http://{self.minikube_ip}:30001/v1/products
   • Check inventory: curl http://{self.minikube_ip}:30002/v1/inventory

✅ Problem Statement 4 - FULLY IMPLEMENTED ON KUBERNETES! ✅
""")

    def run_validation(self):
        """Run complete validation"""
        print("⚓ E-commerce Kubernetes Deployment Validation")
        print("=" * 60)
        print(f"Validation started at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        
        # Get Minikube IP
        if not self.get_minikube_ip():
            return False
            
        # Wait a moment for services to be ready
        time.sleep(2)
        
        # Run validation steps
        steps = [
            self.check_kubernetes_pods,
            self.check_kubernetes_services, 
            self.test_service_health,
            self.test_api_functionality,
            self.check_monitoring_stack
        ]
        
        for step_func in steps:
            if not step_func():
                self.print_error("Validation step failed!")
                return False
                
        # Generate summary
        self.generate_summary_report()
        return True

if __name__ == "__main__":
    validator = KubernetesValidator()
    success = validator.run_validation()
    
    if success:
        print("\n🎉 All validation steps passed! Kubernetes deployment is successful!")
        exit(0)
    else:
        print("\n❌ Validation failed. Please check the deployment.")
        exit(1)