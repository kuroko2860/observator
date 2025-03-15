import { useState } from "react";
import { useSelector } from "react-redux";
import { getAllServices } from "../redux/services/selector";
import axios from "../config/axios";
import PathTree from "../component/PathTree";
import { useNavigate } from "react-router-dom";

const PathSearch = () => {
  const navigate = useNavigate();
  const services = useSelector(getAllServices);

  // State management
  const [operationsByService, setOperationsByService] = useState({});
  const [loadingOperations, setLoadingOperations] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [response, setResponse] = useState(null);
  const [pairs, setPairs] = useState([
    { id: 1, service: "AWS", operation: "" },
  ]);

  const fetchOperations = async (service) => {
    if (operationsByService[service]) return;

    setLoadingOperations((prev) => ({ ...prev, [service]: true }));

    try {
      const { data: operations } = await axios.get(
        `/services/${service}/operations`
      );

      setOperationsByService((prev) => ({
        ...prev,
        [service]: operations,
      }));

      updatePairsWithNewOperation(service, operations);
    } catch (error) {
      console.error(`Error fetching operations for ${service}:`, error);
    } finally {
      setLoadingOperations((prev) => ({ ...prev, [service]: false }));
    }
  };

  const updatePairsWithNewOperation = (service, operations) => {
    setPairs((currentPairs) =>
      currentPairs.map((pair) =>
        pair.service === service &&
        (!pair.operation || !operations.includes(pair.operation))
          ? { ...pair, operation: operations[0] }
          : pair
      )
    );
  };

  const addPair = () => {
    const newId =
      pairs.length > 0 ? Math.max(...pairs.map((p) => p.id)) + 1 : 1;
    const defaultService = services[0];

    setPairs([
      ...pairs,
      {
        id: newId,
        service: defaultService,
        operation: operationsByService[defaultService]?.[0] || "",
      },
    ]);
  };

  const removePair = (id) => {
    setPairs(pairs.filter((pair) => pair.id !== id));
  };

  const updateService = (id, service) => {
    if (!operationsByService[service]) {
      fetchOperations(service);
    }

    setPairs(
      pairs.map((pair) =>
        pair.id === id
          ? {
              ...pair,
              service,
              operation: operationsByService[service]?.[0] || "",
            }
          : pair
      )
    );
  };

  const updateOperation = (id, operation) => {
    setPairs(
      pairs.map((pair) => (pair.id === id ? { ...pair, operation } : pair))
    );
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const { data } = await axios.post("/paths", { pairs });
      setResponse(data);
    } catch (error) {
      console.error("Error submitting data:", error);
      setResponse({
        success: false,
        error: "Failed to submit data",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const allPairsValid = pairs.every(
    (pair) =>
      pair.service && pair.operation && operationsByService[pair.service]
  );

  const renderServiceSelect = (pair) => (
    <div className="w-1/3">
      <label className="block text-sm font-medium mb-1">Service</label>
      <select
        className="w-full rounded-md border border-gray-300 p-2"
        value={pair.service}
        onChange={(e) => updateService(pair.id, e.target.value)}
      >
        {services.map((service) => (
          <option key={service} value={service}>
            {service}
          </option>
        ))}
      </select>
    </div>
  );

  const renderOperationSelect = (pair) => (
    <div className="w-1/3">
      <label className="block text-sm font-medium mb-1">Operation</label>
      <select
        className="w-full rounded-md border border-gray-300 p-2"
        value={pair.operation}
        onChange={(e) => updateOperation(pair.id, e.target.value)}
        disabled={
          loadingOperations[pair.service] || !operationsByService[pair.service]
        }
      >
        {loadingOperations[pair.service] ? (
          <option>Loading operations...</option>
        ) : operationsByService[pair.service] ? (
          operationsByService[pair.service].map((operation) => (
            <option key={operation} value={operation}>
              {operation}
            </option>
          ))
        ) : (
          <option>Select a service first</option>
        )}
      </select>
    </div>
  );

  return (
    <div className="max-w-4xl mx-auto p-6 bg-white rounded-lg shadow-md">
      <h1 className="text-2xl font-bold mb-6">
        Service &amp; Operation Selector
      </h1>

      <form onSubmit={handleSubmit}>
        {pairs.map((pair) => (
          <div key={pair.id} className="flex items-center gap-4 mb-4">
            {renderServiceSelect(pair)}
            {renderOperationSelect(pair)}
            <div className="flex items-end">
              <button
                type="button"
                className="bg-red-500 text-white p-2 rounded-md mt-6"
                onClick={() => removePair(pair.id)}
                disabled={pairs.length <= 1}
              >
                Remove
              </button>
            </div>
          </div>
        ))}

        <div className="flex gap-4 mt-6">
          <button
            type="button"
            className="bg-green-500 text-white py-2 px-4 rounded-md"
            onClick={addPair}
          >
            Add Service & Operation
          </button>

          <button
            type="submit"
            className="bg-blue-500 text-white py-2 px-4 rounded-md"
            disabled={isSubmitting || !allPairsValid}
          >
            {isSubmitting ? "Submitting..." : "Submit"}
          </button>
        </div>
      </form>

      {response && (
        <div className="mt-8 p-4 border rounded-md bg-gray-50">
          <h2 className="text-xl font-bold mb-4">Response</h2>
          <pre className="bg-gray-100 p-4 rounded overflow-x-auto">
            {response.data.paths?.map((path, index) => (
              <div
                key={index}
                onClick={() => navigate(`/path-detail/${path.path_id}`)}
              >
                <h3 className="text-lg font-bold mb-2">Path {index + 1}</h3>
                <PathTree path={path} />
              </div>
            ))}
          </pre>
        </div>
      )}
    </div>
  );
};

export default PathSearch;
