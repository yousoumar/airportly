import { useEffect, useRef, useState } from "react";

interface Data {
    sensorId: number;
    airportId: string;
    sensorType: string;
    value: number;
    timestamp: string;
}
  

const useApp = () => {
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

                const dateRange = await fetchDateInterval(airport, metrics[0]);
                setMetrics(metrics);

                setDateRange(dateRange);
                dateRangeRef.current = dateRange;

                const response = await fetch(
                    `http://localhost:8080/api/v1/${airport}/metric/${isValidMetric(metric, metrics) ? metric : metrics[0]
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

    return {
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
    }
}

export default useApp;