import gatt

manager = gatt.DeviceManager(adapter_name='hci0')

class AnyDevice(gatt.Device):
    def services_resolved(self):
        super().services_resolved()

        # motion
        device_information_service = next(
            s for s in self.services
            if s.uuid == '0000180a-0000-1000-8000-00805f9b34fb')

        firmware_version_characteristic = next(
            c for c in device_information_service.characteristics
            if c.uuid == '00002a29-0000-1000-8000-00805f9b34fb')

        firmware_version_characteristic.enable_notifications()
        firmware_version_characteristic.read_value()


        # motion
        device_information_service = next(
            s for s in self.services
            if s.uuid == '5b2c25e7-7c43-4a15-a4c6-7cf2d81e1b40')

        firmware_version_characteristic = next(
            c for c in device_information_service.characteristics
            if c.uuid == '15e438b8-558e-4b1f-992f-23f90a8c129b')

        firmware_version_characteristic.enable_notifications()
        firmware_version_characteristic.read_value()

        # temperature
        device_information_service = next(
            s for s in self.services
            if s.uuid == 'c355c42e-b56c-458e-bacb-9248717bbac2')

        firmware_version_characteristic = next(
            c for c in device_information_service.characteristics
            if c.uuid == '00002a6e-0000-1000-8000-00805f9b34fb')

        firmware_version_characteristic.enable_notifications()
        firmware_version_characteristic.read_value()

    def characteristic_value_updated(self, characteristic, value):
        print(characteristic.uuid,value[0],value[1])
        # keep reading values if nessecary
        if characteristic.uuid in ['15e438b8-558e-4b1f-992f-23f90a8c129b','00002a6e-0000-1000-8000-00805f9b34fb']:
            characteristic.read_value()

    # def characteristic_enable_notifications_succeeded(self,characteristic,value):
    #     print("HI")

    # def characteristic_enable_notifications_failed(self,characteristic,error):
    #     print("failed enabling notifications: ",characteristic,error)


device = AnyDevice(mac_address='00:0B:57:1B:8C:77', manager=manager)
device.connect()
manager.run()