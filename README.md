# 🏨 Luxe 5-Star Hotel Enterprise Kubernetes Platform

[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.28+-326CE5.svg?logo=kubernetes&logoColor=white)](https://kubernetes.io/)
[![ArgoCD](https://img.shields.io/badge/ArgoCD-GitOps-EF7B4D.svg?logo=argo&logoColor=white)](https://argo-cd.readthedocs.io/)
[![KEDA](https://img.shields.io/badge/KEDA-Event_Driven_Autoscaling-FF6F00.svg)](https://keda.sh/)
[![Kafka](https://img.shields.io/badge/Apache_Kafka-Event_Streams-231F20.svg?logo=apache-kafka&logoColor=white)](https://kafka.apache.org/)
[![AGIC](https://img.shields.io/badge/Azure_AGIC-App_Gateway_Ingress-0078D4.svg)](https://azure.github.io/application-gateway-kubernetes-ingress/)

End-to-end, multi-tier cloud-native application platform powering a **5-Star Luxury Hospitality & Reservation System** on Kubernetes.

---

## 🏛️ Full System Architecture

```
                                [ Public Internet ]
                                         │
                         [ Azure Application Gateway (AGIC) ]
                         (TLS Offloading + WAF + Ingress Controller)
                                         │
        ┌────────────────────────────────┼────────────────────────────────┐
        ▼                                ▼                                ▼
[ / (Frontend Web UI) ]       [ /api/v1 (Booking Engine) ]    [ /events (Ingestion Svc) ]
ClusterIP: Port 80            ClusterIP: Port 8080            ClusterIP: Port 8090
Replicas: 3                   Replicas: 3 (Auto: 3-15 via HPA) Replicas: 2 (KEDA Kafka Scaled)
                                         │                                │
                                         │                                ▼
                                         │                 [ Apache Kafka Event Bus ]
                                         │                 Topic: hotel.reservations.events.v1
                                         ▼
                           [ PaaS Azure Database Service ]
                           (Private Endpoint Integration: 10.101.2.25:5432)
                                         │
                                         ▼
                           [ OpenTelemetry Collector ]
                           (Tracing, Logs, Metrics via ClusterIP: 4317/8889)
```

---

## ⚡ Production Kubernetes Components Included:

| Component | Manifest | Production Role |
|:---|:---|:---|
| **Deployments (3 Replicas)** | `k8s/deployment.yaml` | Non-root, distroless Go backend with rolling updates & anti-affinity. |
| **ClusterIP Service** | `k8s/service.yaml` | Private internal service discovery on Port 8080. |
| **Horizontal Pod Autoscaler** | `k8s/hpa.yaml` | Scales booking engine from 3 to 15 pods based on CPU (70%) and Memory (80%). |
| **KEDA Event Autoscaler** | `k8s/keda-scaledobject.yaml` | Event-driven scaling on Apache Kafka consumer group lag (threshold: 50 msgs). |
| **AGIC Ingress Controller** | `k8s/ingress.yaml` | Azure Application Gateway Ingress with TLS, path-based routing, and WAF. |
| **Frontend Web Service** | `k8s/frontend-service.yaml` | High-availability UI tier on Port 80. |
| **Event Ingestion Worker** | `k8s/ingestion-service.yaml` | High-throughput streaming consumer processing Kafka reservation events. |
| **PaaS Database Link** | `k8s/database-paas-service.yaml` | Azure Database for PostgreSQL linked via Kubernetes ExternalName and Private Endpoints. |
| **Telemetry Collector** | `k8s/telemetry-service.yaml` | OpenTelemetry Collector gathering distributed traces and metrics. |
| **ArgoCD GitOps Sync** | `argocd/hotel-app.yaml` | Automated self-healing GitOps synchronization. |
