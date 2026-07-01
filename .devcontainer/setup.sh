#!/bin/bash
echo "🚀 Installation des outils DevOps manquants..."

# 1. Installation de kind
echo "📦 Installation de kind..."
[ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/latest/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# 2. Installation de ArgoCD CLI
echo "🐙 Installation de ArgoCD CLI..."
curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
sudo install -m 555 argocd-linux-amd64 /usr/local/bin/argocd
rm argocd-linux-amd64

# 3. Installation de Linkerd CLI
echo "🔗 Installation de Linkerd CLI..."
curl --proto '=https' --tlsv1.2 -sSfL https://run.linkerd.io/install | sh
sudo mv ~/.linkerd2/bin/linkerd /usr/local/bin/linkerd

echo "✅ Environnement prêt !"
