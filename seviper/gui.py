import tkinter as tk
import tkinter.ttk as ttk
import tkinter.font as tk_font
import pytkgif
import boa

class ReplayData:
    def __init__(self, file_path):
        data = boa.readlines_txt(file_path, True)
        self.p1_fighter_names = data[-2].split()
        self.p2_fighter_names = data[-1].split()
        keys = ["p1_poke_name", "p1_level", "p1_max_hp", "p1_current_hp",
                "p2_poke_name", "p2_level", "p2_max_hp", "p2_current_hp"]
        funcs = [str, int, int, int, str, int, int, int, str]
        self.battle_surface_statuses = [{key:funcs[i](line.split()[i]) \
                                         for i, key in enumerate(keys)} for line in data[:-3]]

        messages = [" ".join(line.split()[8:]) for line in data[:-3]]
        for i in range(len(self.battle_surface_statuses)):
            self.battle_surface_statuses[i]["message"] = messages[i]

            if self.battle_surface_statuses[i]["p1_poke_name"] == "None":
                self.battle_surface_statuses[i]["p1_poke_name"] = None

            if self.battle_surface_statuses[i]["p1_level"] == -1:
                self.battle_surface_statuses[i]["p1_level"] = None

            if self.battle_surface_statuses[i]["p1_max_hp"] == -1:
                self.battle_surface_statuses[i]["p1_max_hp"] = None

            if self.battle_surface_statuses[i]["p1_current_hp"] == -1:
                self.battle_surface_statuses[i]["p1_current_hp"] = None

            if self.battle_surface_statuses[i]["p2_poke_name"] == "None":
                self.battle_surface_statuses[i]["p2_poke_name"] = None

            if self.battle_surface_statuses[i]["p2_level"] == -1:
                self.battle_surface_statuses[i]["p2_level"] = None

            if self.battle_surface_statuses[i]["p2_max_hp"] == -1:
                self.battle_surface_statuses[i]["p2_max_hp"] = None

            if self.battle_surface_statuses[i]["p2_current_hp"] == -1:
                self.battle_surface_statuses[i]["p2_current_hp"] = None


class HpBar:
    IMAGE_PATH = "C:/Python35/pyckage/seviper/image/hp_bar/"

    def __init__(self, master, max_hp, current_hp, details_font_size):
        assert max_hp >= current_hp

        self.master = master
        self.max_hp = max_hp
        self.current_hp = current_hp

        self.image_label = ttk.Label(master)
        self.details_font = tk_font.Font(family="Lucida Grande", size=details_font_size)
        self.details_label = tk.Label(master=master, text=str(max_hp) + " / " + str(current_hp), font=self.details_font)

        hp_percent = int((self.current_hp / self.max_hp) * 100)
        self.image = tk.PhotoImage(file=HpBar.IMAGE_PATH + str(hp_percent) + ".png")
        self.image_label.configure(image=self.image)
        self.is_animetion_stop = False

    def __update(self):
        hp_percent = int((self.current_hp / self.max_hp) * 100)
        self.image = tk.PhotoImage(file=HpBar.IMAGE_PATH + str(hp_percent) + ".png")
        self.image_label.configure(image=self.image)
        self.details_label.configure(text=str(self.max_hp) + " / " + str(self.current_hp))

    def increment(self):
        self.current_hp += 1
        self.__update()

    def decrement(self):
        self.current_hp -= 1
        self.__update()

    def run_animation(self, next_current_hp, interval):
        if next_current_hp < self.current_hp:
            f = self.decrement
        else:
            f = self.increment

        def run_next_frame_animation(_):
            if self.current_hp == next_current_hp:
                return
            f()
            self.master.after(interval, run_next_frame_animation, None)
        run_next_frame_animation(None)


class Replay:
    def __init__(self, master, file_path):
        master.geometry("610x450")
        self.data = seviper.ReplayData(file_path)
        self.font_size = 15
        font = tk_font.Font(family="Lucida Grande", size=self.font_size)

        self.main_frame = tk.Frame(master=master)

        self.p1_widget_frame = tk.Frame(master=self.main_frame)
        p1_poke_name = self.data.p1_fighter_names[0]
        p1_level = self.data.battle_surface_statuses[0]["p1_level"]

        self.p1_level_label = tk.Label(master=self.p1_widget_frame, text=str(p1_level) + "Lv", font=font)
        self.p1_poke_name_label = tk.Label(master=self.p1_widget_frame, text=p1_poke_name, font=font)
        self.p1_fighter_gifs = {
            poke_name:pytkgif.Gif(self.p1_widget_frame, seviper.POKE_GIF_MIRROR_PATH + poke_name + ".gif") \
            for poke_name in self.data.p1_fighter_names
        }
        self.p1_fighter_gifs[None] = pytkgif.Gif(self.p1_widget_frame, seviper.POKE_GIF_PATH + "None.gif")
        self.p1_hp_bar = HpBar(self.p1_widget_frame, self.data.battle_surface_statuses[0]["p1_max_hp"],
                               self.data.battle_surface_statuses[0]["p1_current_hp"], self.font_size)

        self.p2_widget_frame = tk.Frame(master=self.main_frame)
        p2_poke_name = self.data.p2_fighter_names[0]
        p2_level = self.data.battle_surface_statuses[0]["p2_level"]

        self.p2_level_label = tk.Label(master=self.p2_widget_frame, text=str(p2_level) + "Lv", font=font)
        self.p2_poke_name_label = tk.Label(master=self.p2_widget_frame, text=p2_poke_name, font=font)
        self.p2_fighter_gifs = {
            poke_name:pytkgif.Gif(self.p2_widget_frame, seviper.POKE_GIF_NORMAL_PATH + poke_name + ".gif") \
            for poke_name in self.data.p2_fighter_names
        }
        self.p2_fighter_gifs[None] = pytkgif.Gif(self.p2_widget_frame, seviper.POKE_GIF_PATH + "None.gif")
        self.p2_hp_bar = HpBar(self.p2_widget_frame, self.data.battle_surface_statuses[0]["p2_max_hp"],
                               self.data.battle_surface_statuses[0]["p2_current_hp"], self.font_size)

        self.message_entry = tk.Entry(master=master, width=30, font=font)
        self.master = master

    def run_animation(self):
        level_label_pos = {"row":0, "column":0}
        poke_name_label_pos = {"row":1, "column":0}
        fighter_gif_label_pos = {"row":2, "column":0}
        hp_bar_image_label_pos = {"row":3, "column":0}
        hp_bar_details_label_pos = {"row":4, "column":0}

        self.p1_level_label.grid(**level_label_pos)
        self.p1_poke_name_label.grid(**poke_name_label_pos)
        self.p1_fighter_gifs[self.p1_poke_name_label.cget("text")].image_label.grid(**fighter_gif_label_pos)

        for gif in self.p1_fighter_gifs.values():
            gif.run_animation(50)

        self.p1_hp_bar.image_label.grid(**hp_bar_image_label_pos)
        self.p1_hp_bar.details_label.grid(**hp_bar_details_label_pos)

        self.p2_level_label.grid(**level_label_pos)
        self.p2_poke_name_label.grid(**poke_name_label_pos)
        self.p2_fighter_gifs[self.p2_poke_name_label.cget("text")].image_label.grid(**fighter_gif_label_pos)

        for gif in self.p2_fighter_gifs.values():
            gif.run_animation(50)

        self.p2_hp_bar.image_label.grid(**hp_bar_image_label_pos)
        self.p2_hp_bar.details_label.grid(**hp_bar_details_label_pos)

        self.main_frame.pack()
        self.p1_widget_frame.grid(row=0, column=0)
        self.p2_widget_frame.grid(row=0, column=2)

        row_size, column_size = self.main_frame.grid_size()
        for row in range(row_size):
            self.main_frame.grid_columnconfigure(row, minsize=205)
        self.main_frame.grid_rowconfigure(0, minsize=350)
        self.message_entry.pack()

        interval = 60

        def update_level_label(level_label, next_level):
            str_level = level_label.cget("text")[:2]
            if next_level is None and str_level != "None":
                level_label.grid_forget()
                level_label.configure(text="None")
            elif (str(next_level) != str_level):
                level_label.grid_forget()
                level_label.grid(**level_label_pos)
                level_label.configure(text=str(next_level) + "Lv")

        def update_fighter_gif(fighter_gifs, poke_name_label, next_poke_name):
            poke_name = poke_name_label.cget("text")

            if next_poke_name != poke_name:
                if poke_name == "None":
                    fighter_gifs[None].image_label.grid_forget()
                else:
                    fighter_gifs[poke_name].image_label.grid_forget()

                fighter_gifs[next_poke_name].image_label.grid(**fighter_gif_label_pos)

        def update_poke_name_label(poke_name_label, next_poke_name):
            poke_name = poke_name_label.cget("text")

            if next_poke_name is None:
                poke_name_label.grid_forget()
                poke_name_label.configure(text="None")
            elif next_poke_name != poke_name:
                poke_name_label.grid(**poke_name_label_pos)
                poke_name_label.configure(text=next_poke_name)

        def update_hp_bar(master, hp_bar, next_max_hp, next_current_hp):
            is_hp_bar_animation_phase = False
            if next_max_hp is None:
                hp_bar.image_label.grid_forget()
                hp_bar.details_label.grid_forget()
            elif (next_max_hp != hp_bar.max_hp) or (not hp_bar.image_label.winfo_ismapped()):
                hp_bar.image_label.grid_forget()
                hp_bar.details_label.grid_forget()
                hp_bar = HpBar(master, next_max_hp, next_current_hp, self.font_size)
                hp_bar.image_label.grid(**hp_bar_image_label_pos)
                hp_bar.details_label.grid(**hp_bar_details_label_pos)
            elif next_current_hp != hp_bar.current_hp:
                hp_bar.image_label.grid_forget()
                hp_bar.details_label.grid_forget()
                hp_bar.image_label.grid(**hp_bar_image_label_pos)
                hp_bar.details_label.grid(**hp_bar_details_label_pos)
                if next_current_hp < hp_bar.current_hp:
                    hp_bar.decrement()
                else:
                    hp_bar.increment()
                is_hp_bar_animation_phase = True
            return hp_bar, is_hp_bar_animation_phase

        def run_next_frame_animation(index):
            if len(self.data.battle_surface_statuses) == index:
                return

            battle_surface_scene = self.data.battle_surface_statuses[index]

            next_p1_level = battle_surface_scene["p1_level"]
            next_p2_level = battle_surface_scene["p2_level"]
            next_p1_poke_name = battle_surface_scene["p1_poke_name"]
            next_p2_poke_name = battle_surface_scene["p2_poke_name"]
            next_p1_max_hp = battle_surface_scene["p1_max_hp"]
            next_p2_max_hp = battle_surface_scene["p2_max_hp"]
            next_p1_current_hp = battle_surface_scene["p1_current_hp"]
            next_p2_current_hp = battle_surface_scene["p2_current_hp"]

            update_level_label(self.p1_level_label, next_p1_level)
            update_level_label(self.p2_level_label, next_p2_level)

            update_fighter_gif(self.p1_fighter_gifs, self.p1_poke_name_label, next_p1_poke_name)
            update_fighter_gif(self.p2_fighter_gifs, self.p2_poke_name_label, next_p2_poke_name)

            update_poke_name_label(self.p1_poke_name_label, next_p1_poke_name)
            update_poke_name_label(self.p2_poke_name_label, next_p2_poke_name)

            self.p1_hp_bar, is_p1_hp_bar_animation_phase = \
                update_hp_bar(self.p1_widget_frame, self.p1_hp_bar, next_p1_max_hp, next_p1_current_hp)
            self.p2_hp_bar, is_p2_hp_bar_animation_phase = \
                update_hp_bar(self.p2_widget_frame, self.p2_hp_bar, next_p2_max_hp, next_p2_current_hp)

            next_message = battle_surface_scene["message"]

            if next_message != self.message_entry.cget("text"):
                self.message_entry.delete(0, tk.END)
                self.message_entry.insert(tk.END, next_message)

            if any([is_p1_hp_bar_animation_phase, is_p2_hp_bar_animation_phase]):
                interval = 30
            elif next_message in [seviper.GOOD_EFFECTIVE_BATTLE_MSG,
                                  seviper.BAD_EFFECTIVE_BATTLE_MSG,
                                  seviper.NO_EFFECTIVE_BATTLE_MSG]:
                interval = 600
            elif "戻れ！" in next_message and next_message.count("！") == 2:
                interval = 600
            elif "行け！" in next_message and next_message.count("！") == 2:
                interval = 600
            elif any([move_name in next_message for move_name in seviper.MOVEDEX]) and "！" in next_message:
                interval = 120
            elif ("戻" not in next_message) and ("行" not in next_message) and ("！" in next_message):
                interval = 600
            else:
                interval = 45
            self.master.after(interval, run_next_frame_animation, index + 1)

        run_next_frame_animation(0)

if __name__ == "__main__":
    main_master = tk.Tk()
    main_master.geometry("500x500")

    replay = Replay(main_master, "C:/Python35/pyckage/seviper/replay/0.txt")
    replay_start_button = tk.Button(master=main_master, text="リプレイ開始")
    replay_start_button.pack()

    def replay_start_button_event(event):
        replay_start_button.pack_forget()
        replay.run_animation()

    replay_start_button.bind("<Button-1>", replay_start_button_event)

    main_master.mainloop()
