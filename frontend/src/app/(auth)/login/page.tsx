import { Input } from "@/components/ui/input";


const centerStyle = {
  display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  alignItems: "center",
  minHeight: "100vh", // Adjust height as needed
  maxWidth: "400px", // Adjust maximum width as needed
  margin: "auto", // Center horizontally and vertically
};

const headerStyle = {
  fontSize: "24px",
  fontWeight: "bold",
  marginBottom: "20px",
};

const greenText = {
  color: "#25D366", // WhatsApp green color
};

export default function Login() {
  return (
    <div style={centerStyle}>
      <header style={headerStyle}>
        Login with <span style={greenText}>WhatsApp</span>
      </header>

      <Input
        type="tel"
        placeholder="+11234567890"
        style={{ marginBottom: "10px" }}
      />

      {/* Add small subtitle text that says "Make sure to include the country code" */}
      <small style={{ marginBottom: "20px", color: "gray", fontSize: "12px" }}>
        Make sure to include the country code
      </small>

      <button
        style={{
          backgroundColor: greenText.color,
          color: "white",
          padding: "10px",
          width: "100%",
          border: "none",
          cursor: "pointer",
        }}
      >
        Login
      </button>

    </div>
  );
}