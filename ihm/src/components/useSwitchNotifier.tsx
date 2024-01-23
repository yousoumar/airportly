import { useState, useEffect } from 'react';
import { toast } from "react-toastify";

import mqtt, { MqttClient } from 'mqtt';


function useSwitchNotifier(airport: string) {
    const [isSwitchOn, setIsSwitchOn] = useState(false);
    const [mqttClient, setMqttClient] = useState<MqttClient | null>(null);


    const handleSwitchToggle = () => {
        if (isSwitchOn) {
            if (mqttClient) {
                mqttClient.end();
                setMqttClient(null);
            }
        } else {
            const newMqttClient = mqtt.connect('ws://localhost:9001/');
            console.log("Connected");
            
            newMqttClient.on('connect', () => {
                newMqttClient?.subscribe(`airport/${airport}/alert/#`);
            });

            newMqttClient.on('message', (_topic, message) => {
                const newAlert = JSON.parse(message.toString());                
                toast.warn(newAlert?.Message)
            });

            setMqttClient(newMqttClient);
        }

        setIsSwitchOn(!isSwitchOn);
    };

    useEffect(() => {
        return () => {
            if (mqttClient) {
                mqttClient.end();
            }
        };
    }, [mqttClient]);

    return { isSwitchOn, handleSwitchToggle }

}

export default useSwitchNotifier