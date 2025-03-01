import { useState } from "react";
import { Clock, Calendar, ArrowRight, Check } from "lucide-react";

const CheckInOutSelector = () => {
  const [checkInDateTime, setCheckInDateTime] = useState("");
  const [checkOutDateTime, setCheckOutDateTime] = useState("");
  const [validationError, setValidationError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);
  const [activeField, setActiveField] = useState(null);

  // Get current date and time in the format YYYY-MM-DDThh:mm
  const getCurrentDateTime = () => {
    const now = new Date();
    return now.toISOString().slice(0, 16);
  };

  // Calculate minimum checkout time (must be at least 1 hour after check-in)
  const getMinCheckoutTime = () => {
    if (!checkInDateTime) return getCurrentDateTime();

    const checkInDate = new Date(checkInDateTime);
    checkInDate.setHours(checkInDate.getHours() + 1);
    return checkInDate.toISOString().slice(0, 16);
  };

  const handleCheckInChange = (e) => {
    const newCheckInDateTime = e.target.value;
    setCheckInDateTime(newCheckInDateTime);
    setValidationError("");

    if (checkOutDateTime) {
      const checkIn = new Date(newCheckInDateTime);
      const checkOut = new Date(checkOutDateTime);

      if (checkOut <= checkIn) {
        const updatedCheckOut = new Date(checkIn);
        updatedCheckOut.setHours(updatedCheckOut.getHours() + 1);
        setCheckOutDateTime(updatedCheckOut.toISOString().slice(0, 16));
      }
    }
  };

  const handleCheckOutChange = (e) => {
    const newCheckOutDateTime = e.target.value;
    setCheckOutDateTime(newCheckOutDateTime);

    if (checkInDateTime && newCheckOutDateTime) {
      const checkIn = new Date(checkInDateTime);
      const checkOut = new Date(newCheckOutDateTime);

      if (checkOut <= checkIn) {
        setValidationError("Check-out must be after check-in");
      } else {
        setValidationError("");
      }
    }
  };

  const formatDateDisplay = (dateString) => {
    if (!dateString) return "";
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      weekday: "short",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const calculateDuration = () => {
    if (!checkInDateTime || !checkOutDateTime) return null;

    const checkIn = new Date(checkInDateTime);
    const checkOut = new Date(checkOutDateTime);
    const diff = checkOut - checkIn;

    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));

    let result = "";
    if (days > 0) result += `${days} day${days > 1 ? "s" : ""}`;
    if (hours > 0)
      result += `${days > 0 ? " " : ""}${hours} hour${hours > 1 ? "s" : ""}`;

    return result;
  };

  const handleSubmit = () => {
    if (!checkInDateTime || !checkOutDateTime) {
      setValidationError("Both check-in and check-out times are required");
      return;
    }

    const checkIn = new Date(checkInDateTime);
    const checkOut = new Date(checkOutDateTime);

    if (checkOut <= checkIn) {
      setValidationError("Check-out must be after check-in");
      return;
    }

    setValidationError("");
    setIsSubmitting(true);

    // Simulate API call
    setTimeout(() => {
      setIsSubmitting(false);
      setIsSuccess(true);

      // Reset success message after 3 seconds
      setTimeout(() => {
        setIsSuccess(false);
      }, 3000);
    }, 1500);
  };

  return (
    <div className="flex flex-col items-center justify-center w-full p-6 bg-gradient-to-br from-gray-900 to-gray-800 text-white rounded-lg shadow-xl">
      <h2 className="text-2xl font-bold mb-6 text-center">Book Your Stay</h2>

      <div className="flex flex-col sm:flex-row w-full max-w-3xl gap-4 mb-6">
        <div className="relative w-full sm:w-1/2 group">
          <div
            className={`absolute -inset-0.5 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg blur opacity-30 group-hover:opacity-100 transition duration-300 ${
              activeField === "check-in" ? "opacity-75" : ""
            }`}
          ></div>
          <div className="relative bg-gray-800 p-4 rounded-lg">
            <label
              htmlFor="check-in"
              className="flex items-center text-lg font-medium mb-2"
            >
              <Calendar size={18} className="mr-2 text-blue-400" />
              Check-in
            </label>
            <div className="relative">
              <input
                id="check-in"
                type="datetime-local"
                value={checkInDateTime}
                onChange={handleCheckInChange}
                min={getCurrentDateTime()}
                onFocus={() => setActiveField("check-in")}
                onBlur={() => setActiveField(null)}
                className="w-full p-3 bg-gray-700 border border-gray-600 rounded-md text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all duration-300"
              />
              {checkInDateTime && (
                <div className="mt-2 text-blue-300 font-medium">
                  {formatDateDisplay(checkInDateTime)}
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="flex items-center justify-center my-2 sm:my-0">
          <div className="w-10 h-10 flex items-center justify-center bg-blue-600 rounded-full shadow-lg">
            <ArrowRight size={18} />
          </div>
        </div>

        <div className="relative w-full sm:w-1/2 group">
          <div
            className={`absolute -inset-0.5 bg-gradient-to-r from-purple-600 to-pink-500 rounded-lg blur opacity-30 group-hover:opacity-100 transition duration-300 ${
              activeField === "check-out" ? "opacity-75" : ""
            }`}
          ></div>
          <div className="relative bg-gray-800 p-4 rounded-lg">
            <label
              htmlFor="check-out"
              className="flex items-center text-lg font-medium mb-2"
            >
              <Calendar size={18} className="mr-2 text-purple-400" />
              Check-out
            </label>
            <div className="relative">
              <input
                id="check-out"
                type="datetime-local"
                value={checkOutDateTime}
                onChange={handleCheckOutChange}
                min={getMinCheckoutTime()}
                onFocus={() => setActiveField("check-out")}
                onBlur={() => setActiveField(null)}
                className={`w-full p-3 bg-gray-700 border rounded-md text-white focus:ring-2 focus:ring-purple-500 focus:border-transparent outline-none transition-all duration-300 ${
                  validationError && validationError.includes("Check-out")
                    ? "border-red-500"
                    : "border-gray-600"
                }`}
              />
              {checkOutDateTime && (
                <div className="mt-2 text-purple-300 font-medium">
                  {formatDateDisplay(checkOutDateTime)}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {checkInDateTime && checkOutDateTime && !validationError && (
        <div className="w-full max-w-3xl mb-6">
          <div className="bg-gray-800/50 rounded-lg p-4 border border-gray-700">
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <Clock size={18} className="mr-2 text-green-400" />
                <span className="text-gray-300">Duration:</span>
              </div>
              <span className="font-bold text-white">
                {calculateDuration()}
              </span>
            </div>
          </div>
        </div>
      )}

      {validationError && (
        <div className="w-full max-w-3xl mb-6">
          <div className="bg-red-900/20 text-red-400 p-3 rounded-md border border-red-800/50 flex items-start">
            <div className="mr-2 mt-0.5">⚠️</div>
            <div>{validationError}</div>
          </div>
        </div>
      )}

      <button
        onClick={handleSubmit}
        disabled={
          isSubmitting ||
          isSuccess ||
          !!validationError ||
          !checkInDateTime ||
          !checkOutDateTime
        }
        className={`
          relative overflow-hidden w-full max-w-3xl py-3 px-6 rounded-lg font-medium text-lg 
          transition-all duration-300 shadow-lg
          ${
            isSubmitting ||
            isSuccess ||
            !!validationError ||
            !checkInDateTime ||
            !checkOutDateTime
              ? "bg-gray-600 cursor-not-allowed opacity-70"
              : "bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 cursor-pointer"
          }
        `}
      >
        <div className="relative z-10 flex items-center justify-center">
          {isSubmitting ? (
            <div className="flex items-center">
              <div className="animate-spin h-5 w-5 mr-2 border-2 border-white border-t-transparent rounded-full"></div>
              Processing...
            </div>
          ) : isSuccess ? (
            <div className="flex items-center text-green-300">
              <Check size={18} className="mr-2" />
              Booking Confirmed!
            </div>
          ) : (
            "Confirm Reservation"
          )}
        </div>

        {/* Button hover effect */}
        <div className="absolute inset-0 h-full w-full scale-0 rounded-lg transition-all duration-300 group-hover:scale-100 group-hover:bg-white/10"></div>
      </button>
    </div>
  );
};

export default CheckInOutSelector;
