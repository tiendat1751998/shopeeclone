"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useCartStore } from "@/lib/store/cart";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Price } from "@/components/ui/Price";
import { ordersApi } from "@/lib/api/orders";
import { ApiError } from "@/lib/api/client";

function sanitizeAddressValue(val: string): string {
  return val.trim().slice(0, 500).replace(/[<>]/g, "");
}

export default function CheckoutPage() {
  const router = useRouter();
  const { selectedItems, subtotal, clearCart } = useCartStore();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [address, setAddress] = useState({ name: "", phone: "", address_line1: "", city: "", state: "", postal_code: "", country: "SG" });
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [paymentMethod, setPaymentMethod] = useState("shopeepay");

  const validate = (): boolean => {
    const errors: Record<string, string> = {};
    if (!address.name.trim()) errors.name = "Name is required";
    if (!address.phone.trim()) errors.phone = "Phone is required";
    else if (!/^[+\d\s()-]{7,20}$/.test(address.phone.trim())) errors.phone = "Invalid phone number";
    if (!address.address_line1.trim()) errors.address_line1 = "Address is required";
    if (!address.city.trim()) errors.city = "City is required";
    if (!address.postal_code.trim()) errors.postal_code = "Postal code is required";
    setFieldErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async () => {
    if (selectedItems().length === 0) return;
    if (!validate()) return;
    setIsSubmitting(true);
    setError("");
    try {
      const items = selectedItems();
      const firstItem = items[0];
      await ordersApi.checkout({
        items: items.map((i) => ({
          product_id: i.product_id,
          sku_id: i.sku_id,
          shop_id: i.shop_id,
          name: i.name,
          quantity: i.quantity,
          unit_price: i.price,
          image_url: i.image_url,
        })),
        seller_id: firstItem?.shop_id || "",
        shipping_address: {
          street1: sanitizeAddressValue(address.address_line1),
          city: sanitizeAddressValue(address.city),
          state: sanitizeAddressValue(address.state),
          postal_code: sanitizeAddressValue(address.postal_code),
          country: "SG",
          phone: sanitizeAddressValue(address.phone),
        },
        billing_address: {
          street1: sanitizeAddressValue(address.address_line1),
          city: sanitizeAddressValue(address.city),
          state: sanitizeAddressValue(address.state),
          postal_code: sanitizeAddressValue(address.postal_code),
          country: "SG",
          phone: sanitizeAddressValue(address.phone),
        },
      });
      clearCart();
      router.push("/account");
    } catch (e: unknown) {
      if (e instanceof ApiError && e.status === 429) {
        setError("Too many requests. Please wait a moment and try again.");
      } else {
        setError(e instanceof Error ? e.message : "Checkout failed. Please try again.");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  if (selectedItems().length === 0) return <div className="container py-16 text-center"><h2 className="text-xl font-semibold mb-2">No items to checkout</h2><Button variant="primary" onClick={() => router.push("/cart")}>Go to Cart</Button></div>;

  return (
    <div className="container py-6">
      <h1 className="text-xl font-semibold mb-6">Checkout</h1>
      {error && <div className="bg-red-50 text-red-600 text-sm p-3 rounded mb-4">{error}</div>}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <div className="card p-6">
            <h3 className="font-semibold mb-4">Shipping Address</h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <Input label="Full Name" value={address.name} onChange={(e) => { setAddress({ ...address, name: e.target.value }); setFieldErrors((p) => ({ ...p, name: "" })); }} error={fieldErrors.name} required maxLength={100} />
              <Input label="Phone" value={address.phone} onChange={(e) => { setAddress({ ...address, phone: e.target.value }); setFieldErrors((p) => ({ ...p, phone: "" })); }} error={fieldErrors.phone} required maxLength={20} />
              <Input label="Address" value={address.address_line1} onChange={(e) => { setAddress({ ...address, address_line1: e.target.value }); setFieldErrors((p) => ({ ...p, address_line1: "" })); }} error={fieldErrors.address_line1} required className="sm:col-span-2" maxLength={500} />
              <Input label="City" value={address.city} onChange={(e) => { setAddress({ ...address, city: e.target.value }); setFieldErrors((p) => ({ ...p, city: "" })); }} error={fieldErrors.city} required maxLength={100} />
              <Input label="Postal Code" value={address.postal_code} onChange={(e) => { setAddress({ ...address, postal_code: e.target.value }); setFieldErrors((p) => ({ ...p, postal_code: "" })); }} error={fieldErrors.postal_code} required maxLength={20} />
            </div>
          </div>
          <div className="card p-6">
            <h3 className="font-semibold mb-4">Payment Method</h3>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
              {(["shopeepay", "credit_card", "paynow", "grabpay"] as const).map((m) => (
                <button key={m} type="button" onClick={() => setPaymentMethod(m)}
                  className={`p-3 border rounded text-sm font-medium transition-colors ${paymentMethod === m ? "border-[#ee4d2d] bg-[#fff0ed] text-[#ee4d2d]" : "border-[#e8e8e8] hover:border-[#ee4d2d]"}`}>
                  {m === "shopeepay" ? "ShopeePay" : m === "credit_card" ? "Credit Card" : m === "paynow" ? "PayNow" : "GrabPay"}
                </button>
              ))}
            </div>
          </div>
        </div>
        <div className="card p-4 h-fit sticky top-24">
          <h3 className="font-semibold mb-4">Order Summary</h3>
          <div className="border-t border-[#e8e8e8] pt-3 space-y-2 text-sm">
            <div className="flex justify-between"><span className="text-[#757575]">Subtotal</span><span>S${subtotal().toFixed(2)}</span></div>
            <div className="flex justify-between"><span className="text-[#757575]">Shipping</span><span className="text-[#00bfa5]">Free</span></div>
            <div className="flex justify-between font-semibold pt-2 border-t border-[#e8e8e8]"><span>Total</span><Price amount={subtotal()} size="md" /></div>
          </div>
          <Button variant="primary" fullWidth isLoading={isSubmitting} onClick={handleSubmit} className="mt-4">Place Order</Button>
        </div>
      </div>
    </div>
  );
}
