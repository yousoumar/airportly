/* eslint-disable @typescript-eslint/no-explicit-any */
import { Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent } from "@mui/material";
import { LineChart } from "@mui/x-charts";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";

import dayjs from "dayjs";
import { useEffect, useRef, useState } from "react";
import "./App.css";
export interface Data {
  sensorId: number;
  airportId: string;
  sensorType: string;
  value: number;
  timestamp: string;
}

function App() {
  const [data, setData] = useState<Data[]>([]);
  const [loadingData, setLoadingData] = useState(true);

  const [width, setWidth] = useState(window.innerWidth);
  const [metric, setMetric] = useState("");
  const [metrics, setMetrics] = useState<string[]>([]);
  const [airport, setAirport] = useState("");
  const [error, setError] = useState("");
  const dateRangeRef = useRef<{ startTime: string; endTime: string }>({
    startTime: "",
    endTime: "",
  });
  const [dateRange, setDateRange] = useState<{ startTime: string; endTime: string }>({
    startTime: "",
    endTime: "",
  });

  const [airports, setAirports] = useState<string[] | null>(null);
  useEffect(() => {
    const handleResize = () => setWidth(window.innerWidth);
    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  useEffect(() => {
    const fetchAvailableAirports = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/v1/metadata/airports`);
        const data = await response.json();
        setAirports(data);
        setAirport(data[0]);
        return data[0];
      } catch (error: any) {
        setError(error.message);
      }
    };
    fetchAvailableAirports();
  }, []);

  useEffect(() => {
    if (!airport) {
      return;
    }

    const fetchAvailableMetrics = async (airport: string) => {
      try {
        const response = await fetch(`http://localhost:8080/api/v1/${airport}/available-metrics`);
        const data = await response.json();
        return data;
      } catch (error: any) {
        setError(error.message);
        console.error(error);
        return null;
      }
    };

    const fetchDateInterval = async (airport: string, metric: string) => {
      try {
        const response = await fetch(
          `http://localhost:8080/api/v1/${airport}/metric/${metric}/date-range`
        );
        const data = await response.json();
        return data;
      } catch (error: any) {
        setError(error.message);
        return null;
      }
    };

    const isValidMetric = (metric: string, metrics: string[]) => metric && metrics.includes(metric);
    const fetchData = async () => {
      try {
        setLoadingData(true);
        const metrics = await fetchAvailableMetrics(airport);

        if (!isValidMetric(metric, metrics)) {
          setMetric(metrics[0]);
        }

        console.log(metrics);
        const dateRange = await fetchDateInterval(airport, metrics[0]);
        console.log(dateRange);
        setMetrics(metrics);

        setDateRange(dateRange);
        console.log(dateRange);
        dateRangeRef.current = dateRange;

        const response = await fetch(
          `http://localhost:8080/api/v1/${airport}/metric/${
            isValidMetric(metric, metrics) ? metric : metrics[0]
          }?startTime=${dateRange.startTime}&endTime=${dateRange.endTime}`
        );
        const data = await response.json();
        setData(data);
      } catch (error: any) {
        setError(error.message);
        console.error(error);
      } finally {
        setLoadingData(false);
      }
    };

    fetchData();
  }, [airport, metric]);

  const fetchDataOnDateChange = async (
    airport: string,
    metric: string,
    startTime: string,
    endTime: string
  ) => {
    try {
      setLoadingData(true);
      const response = await fetch(
        `http://localhost:8080/api/v1/${airport}/metric/${metric}?startTime=${startTime}&endTime=${endTime}`
      );
      const data = await response.json();
      setData(data);
    } catch (error: any) {
      setError(error.message);
      console.error(error);
    } finally {
      setLoadingData(false);
    }
  };

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
    </main>
  );
}

export default App;
