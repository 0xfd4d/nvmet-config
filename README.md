### nvmet-config
Cli tool for configuring Nvme-oF Linux target

### Example
This will apply content of file to target:

`# nvmet-config import nvmet.yaml`

###### nvmet.yaml:
```
ports:
  - name: 1
    addr_adrfam: ipv4
    addr_traddr: 10.0.0.1
    addr_trsvcid: 4420
    addr_trtype: tcp
    subsystems:
      - name: example
subsystems:
  - name: example
    namespaces:
      - name: 1
        enable: 1
        device_path: /dev/test1
        device_uuid: c37dc686-9999-45f4-b7ae-6f49b20ed558
      - name: 2
        enable: 1
        device_path: /dev/test2
        device_uuid: aac800fd-aca9-4a50-9395-67d41ccbe508
```
