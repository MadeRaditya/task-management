import { useEffect, useState } from "react";
import {
  ActivityIndicator,
  Modal,
  Pressable,
  StyleSheet,
  Text,
  TextInput,
  View,
} from "react-native";
import { Task, TaskStatus } from "../types/task";

interface Props {
  visible: boolean;
  task: Task | null;
  onClose: () => void;
  onSubmit: (
    id: number,
    payload: {
      title: string;
      description: string;
      status: TaskStatus;
      assignee: string;
    },
  ) => Promise<void>;
}

const STATUSES: TaskStatus[] = ["todo", "in_progress", "done"];

export function EditTaskModal({ visible, task, onClose, onSubmit }: Props) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [assignee, setAssignee] = useState("");
  const [status, setStatus] = useState<TaskStatus>("todo");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (task) {
      setTitle(task.title);
      setDescription(task.description ?? "");
      setAssignee(task.assignee ?? "");
      setStatus(task.status);
      setError(null);
    }
  }, [task]);

  const handleSave = async () => {
    if (!task) return;
    if (!title.trim()) {
      setError("Title wajib diisi");
      return;
    }
    setSaving(true);
    setError(null);
    try {
      await onSubmit(task.id, {
        title: title.trim(),
        description: description.trim(),
        status,
        assignee: assignee.trim(),
      });
      onClose();
    } catch (e: any) {
      setError(e?.message ?? "Gagal menyimpan");
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={onClose}
    >
      <View style={styles.backdrop}>
        <View style={styles.card} testID="edit-modal">
          <Text style={styles.title}>Edit Task</Text>

          <Text style={styles.label}>Title</Text>
          <TextInput
            testID="edit-title"
            style={styles.input}
            value={title}
            onChangeText={setTitle}
          />

          <Text style={styles.label}>Description</Text>
          <TextInput
            testID="edit-description"
            style={[styles.input, styles.multiline]}
            value={description}
            onChangeText={setDescription}
            multiline
          />

          <Text style={styles.label}>Assignee</Text>
          <TextInput
            testID="edit-assignee"
            style={styles.input}
            value={assignee}
            onChangeText={setAssignee}
          />

          <Text style={styles.label}>Status</Text>
          <View style={styles.statusRow}>
            {STATUSES.map((s) => {
              const active = status === s;
              return (
                <Pressable
                  key={s}
                  testID={`edit-status-${s}`}
                  onPress={() => setStatus(s)}
                  style={[styles.chip, active && styles.chipActive]}
                >
                  <Text
                    style={[styles.chipText, active && styles.chipTextActive]}
                  >
                    {s}
                  </Text>
                </Pressable>
              );
            })}
          </View>

          {error && <Text style={styles.error}>{error}</Text>}

          <View style={styles.actions}>
            <Pressable
              testID="edit-cancel"
              onPress={onClose}
              style={[styles.btn, styles.btnGhost]}
              disabled={saving}
            >
              <Text style={styles.btnGhostText}>Batal</Text>
            </Pressable>
            <Pressable
              testID="edit-save"
              onPress={handleSave}
              style={[styles.btn, styles.btnPrimary]}
              disabled={saving}
            >
              {saving ? (
                <ActivityIndicator color="#fff" />
              ) : (
                <Text style={styles.btnPrimaryText}>Simpan</Text>
              )}
            </Pressable>
          </View>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  backdrop: {
    flex: 1,
    backgroundColor: "rgba(0,0,0,0.4)",
    alignItems: "center",
    justifyContent: "center",
    padding: 16,
  },
  card: {
    width: "100%",
    maxWidth: 480,
    backgroundColor: "#fff",
    borderRadius: 12,
    padding: 16,
  },
  title: {
    fontSize: 18,
    fontWeight: "700",
    marginBottom: 12,
    color: "#111827",
  },
  label: { fontSize: 13, color: "#374151", marginTop: 8, marginBottom: 4 },
  input: {
    borderWidth: 1,
    borderColor: "#d1d5db",
    borderRadius: 8,
    paddingHorizontal: 12,
    paddingVertical: 8,
    fontSize: 15,
  },
  multiline: { minHeight: 70, textAlignVertical: "top" },
  statusRow: { flexDirection: "row", gap: 8, marginTop: 4 },
  chip: {
    paddingHorizontal: 10,
    paddingVertical: 6,
    borderRadius: 999,
    borderWidth: 1,
    borderColor: "#d1d5db",
  },
  chipActive: { backgroundColor: "#2563eb", borderColor: "#2563eb" },
  chipText: { color: "#374151", fontSize: 12 },
  chipTextActive: { color: "#fff", fontWeight: "600" },
  error: { color: "#dc2626", marginTop: 8 },
  actions: {
    flexDirection: "row",
    justifyContent: "flex-end",
    gap: 8,
    marginTop: 16,
  },
  btn: {
    paddingHorizontal: 16,
    paddingVertical: 10,
    borderRadius: 8,
    minWidth: 90,
    alignItems: "center",
  },
  btnGhost: { backgroundColor: "#f3f4f6" },
  btnGhostText: { color: "#374151" },
  btnPrimary: { backgroundColor: "#2563eb" },
  btnPrimaryText: { color: "#fff", fontWeight: "600" },
});
