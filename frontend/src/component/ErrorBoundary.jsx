import * as React from "react";

const ErrorBoundary = ({ fallback, children }) => {
  const [hasError, setHasError] = React.useState(false);

  const getDerivedStateFromError = (error) => {
    // Update state so the next render will show the fallback UI.
    return { hasError: true };
  };

  const componentDidCatch = (error, errorInfo) => {
    logErrorToMyService(
      error,
      // Example "componentStack":
      //   in ComponentThatThrows (created by App)
      //   in ErrorBoundary (created by App)
      //   in div (created by App)
      //   in App
      errorInfo.componentStack,
      // Only available in react@canary.
      // Warning: Owner Stack is not available in production.
      React.captureOwnerStack()
    );
  };
  const logErrorToMyService = (error, componentStack, ownerStack) => {
    // Implementation for logging error to your service
    console.error("Error logged:", error, componentStack, ownerStack);
    // You can send this information to an external logging service here
  };

  if (hasError) {
    // You can render any custom fallback UI
    return fallback;
  }

  return children;
};

export default ErrorBoundary;
