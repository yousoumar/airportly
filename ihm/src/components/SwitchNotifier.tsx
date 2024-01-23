import useSwitchNotifier from "./useSwitchNotifier";
import 'react-toastify/dist/ReactToastify.css';
import Switch from "@mui/material/Switch/Switch";

interface Props {
    airport: string;
}

function SwitchNotifier(props: Props) {
    const { isSwitchOn, handleSwitchToggle } = useSwitchNotifier(props.airport)
    return (
        <div>
            <label>
                Alerts :    Off
                <Switch
                    checked={isSwitchOn}
                    onChange={handleSwitchToggle}
                    inputProps={{ 'aria-label': 'controlled' }}
                /> On
            </label>
        </div>

    )
}

export default SwitchNotifier