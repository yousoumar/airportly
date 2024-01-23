/* eslint-disable @typescript-eslint/no-explicit-any */
import { Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent } from "@mui/material";
import { LineChart } from "@mui/x-charts";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";
import { Slide, ToastContainer } from "react-toastify";

import useApp from "./useApp";

import dayjs from "dayjs";
import "./App.css";
import SwitchNotifier from "./components/SwitchNotifier";


function App() {
  const {
    loadingData,
    airport,
    airports,
    data,
    metric,
    metrics,
    dateRange,
    dateRangeRef,
    error,
    width,
    setAirport,
    setMetric,
    setDateRange,
    fetchDataOnDateChange
  } = useApp();

  if (loadingData) {
    return (
      <Box
        style={{
          display: "flex",
          justifyContent: "center",
          height: "100vh",
          paddingTop: "7rem",
        }}
      >
        Loading...
      </Box>
    );
  }

  return (
    <main style={{ marginTop: "4rem" }}>
      <Box
        sx={{
          display: "flex",
          gap: "1rem",
          alignItems: "center",
          marginLeft: "2.5rem",
        }}
      >
        <Box sx={{ minWidth: 120 }}>
          <FormControl fullWidth>
            <InputLabel>Airport</InputLabel>
            <Select
              value={airport}
              label="Metric"
              onChange={(event: SelectChangeEvent) => {
                setAirport(event.target.value as string);
              }}
            >
              {airports?.map((airport, index) => (
                <MenuItem key={index} value={airport}>
                  {airport}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Box>
        <Box sx={{ minWidth: 120 }}>
          <FormControl fullWidth>
            <InputLabel>Metric</InputLabel>
            <Select
              value={metric}
              label="Metric"
              defaultValue={metrics[0]}
              onChange={(event: SelectChangeEvent) => {
                setMetric(event.target.value as string);
              }}
            >
              {metrics.map((metric, index) => (
                <MenuItem key={index} value={metric}>
                  {metric}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Box>

        <DateTimePicker
          label="Start date"
          minDate={dayjs(dateRangeRef.current?.startTime)}
          maxDate={dayjs(dateRangeRef.current?.endTime)}
          value={dayjs(dateRange?.startTime)}
          defaultValue={dayjs(dateRangeRef.current?.startTime)}
          onAccept={(date: dayjs.Dayjs | null) => {
            setDateRange((prev) => ({ ...prev, startTime: date!.toISOString() }));
            fetchDataOnDateChange(airport, metric, date!.toISOString(), dateRange.endTime);
          }}
        />
        <DateTimePicker
          value={dayjs(dateRange?.endTime)}
          minDate={dayjs(dateRangeRef.current?.startTime)}
          maxDate={dayjs(dateRangeRef.current?.endTime)}
          label="End date"
          onAccept={(date: dayjs.Dayjs | null) => {
            setDateRange((prev) => ({ ...prev, endTime: date!.toISOString() ?? "" }));
            fetchDataOnDateChange(airport, metric, dateRange.startTime, date!.toISOString());
          }}
          defaultValue={dayjs(dateRangeRef.current?.startTime)}
        />
        <Box sx={{ border: 0.5, borderColor: 'grey.400', minWidth: 120, padding:1, borderRadius: 1}}>
          <SwitchNotifier airport={airport} />
        </Box>

      </Box>
      {!data?.length && !error ? (
        <div style={{ marginTop: "6rem" }}>
          No data for your airport, metric or time interval selection :(
        </div>
      ) : null}
      {data?.length && !error ? (
        <Box>
          <LineChart
            xAxis={[
              {
                id: "timestamp",
                data: data.map((d) => new Date(d.timestamp).toLocaleTimeString()),
                scaleType: "band",
              },
            ]}
            series={[
              {
                data: data.map((d) => d.value),
              },
            ]}
            width={width}
            height={300}
          />
        </Box>
      ) : null}
      {error && <div style={{ marginTop: "6rem" }}>{error}</div>}
      <ToastContainer
        position="top-right"
        autoClose={5000}
        hideProgressBar={false}
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
        theme="dark"
        transition={Slide}
      />
    </main>
  );
}

export default App;
