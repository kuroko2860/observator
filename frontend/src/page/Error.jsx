function Error({ error }) {
  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        minHeight: "100vh",
        backgroundColor: "#f7f7f7",
        color: "#333",
        fontSize: "2rem",
      }}
    >
      Error: {error.message}
    </div>
  );
}

export default Error;
