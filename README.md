# Time Scaler Operator

A simple Kubernetes operator made using the `operator-sdk` framework.

This operator allows you to define a `TimeScale` custom resource that enables time-based scaling of your deployments.
```
apiVersion: api.engineerkenji.com/v1alpha1
kind: TimeScaler
metadata:
  name: timescaler-mysql
spec:
  startTime: 0           // 12:00 AM UTC
  endTime: 10            // 10:00 AM UTC
  replicas: 3            // Update deployment to 3 replicas
  deployments:
    - name: mysql        // deployment name to scale
      namespace: default // namespace of deployment

```
