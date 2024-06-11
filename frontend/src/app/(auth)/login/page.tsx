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

export default function Login() {
  return (
    <div style={centerStyle}>
      <header style={headerStyle}>
        {/* Login with email, looking pretty with a nice font and email icon */}
        Login with email  📧
      </header>

      <Input
        type="email"
        placeholder="email"
        style={{ marginBottom: "10px" }}
      />
    </div>
  );
}
