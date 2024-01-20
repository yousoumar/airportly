/* eslint-disable @typescript-eslint/no-explicit-any */
import { Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent } from "@mui/material";
import { LineChart } from "@mui/x-charts";
import { useEffect, useState } from "react";
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
  const [loadingMetrics, setLoadingMetrics] = useState(true);
  const [width, setWidth] = useState(window.innerWidth);
  const [metric, setMetric] = useState("");
  const [metrics, setMetrics] = useState<string[]>([]);
  const [airport, setAirport] = useState("");
  const [error, setError] = useState("");

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
        console.log(data);
        setAirports(data);
        setAirport(data[0]);
        return data[0];
      } catch (error: any) {
        setError(error.message);
      }
    };

    const fetchAvailableMetrics = async (airport: string) => {
      try {
        setLoadingMetrics(true);
        const response = await fetch(`http://localhost:8080/api/v1/${airport}/available-metrics`);
        const data = await response.json();
        console.log(data);
        setMetrics(data);
        setMetric(data[0]);
      } catch (error: any) {
        setError(error.message);
      } finally {
        setLoadingMetrics(false);
      }
    };
    if (!airports) {
      fetchAvailableAirports();
    }
    if (airport) {
      fetchAvailableMetrics(airport);
    }
  }, [airport, airports]);

  useEffect(() => {
    if (!airport || !metric || loadingMetrics) {
      return;
    }
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

    const fetchData = async () => {
      try {
        setLoadingData(true);
        const dateRange = await fetchDateInterval(airport, metric);

        const response = await fetch(
          `http://localhost:8080/api/v1/${airport}/metric/${metric}?startTime=${dateRange.startTime}&endTime=${dateRange.endTime}`
        );
        const data = await response.json();

        setData(data);
      } catch (error: any) {
        setError(error.message);
      } finally {
        setLoadingData(false);
      }
    };

    fetchData();
  }, [metric, loadingMetrics]);

  if (loadingData || loadingMetrics) {
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
      </Box>
      {!data?.length && !error ? (
        <div>No data for your airport, metric or time interval selection :(</div>
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
      {error && <div style={{ marginTop: "2rem" }}>{error}</div>}
    </main>
  );
}

export default App;
