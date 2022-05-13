package doctorday

import (
	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/EventStore/training-introduction-go/infrastructure"
)

type CommandHandlers struct {
	*infrastructure.CommandHandlerBase
}

func NewHandlers(repository DayRepository) CommandHandlers {
	commandHandler := CommandHandlers{infrastructure.NewCommandHandler()}

	commandHandler.Register(commands.ScheduleDay{}, func(c infrastructure.Command, m infrastructure.CommandMetadata) error {
		cmd := c.(commands.ScheduleDay)
		id := NewDayID(NewDoctorID(cmd.DoctorId), cmd.Date)
		day, err := repository.Get(id)
		if err != nil {
			return err
		}

		err = day.ScheduleDay(NewDoctorID(cmd.DoctorId), cmd.Date, cmd.Slots)
		if err != nil {
			return err
		}

		repository.Save(day, m)
		return nil
	})

	commandHandler.Register(commands.BookSlot{}, func(c infrastructure.Command, m infrastructure.CommandMetadata) error {
		cmd := c.(commands.BookSlot)
		day, err := repository.Get(NewDayIDFrom(cmd.DayId))
		if err != nil {
			return err
		}

		err = day.BookSlot(NewSlotID(cmd.SlotId), NewPatientID(cmd.PatientId))
		if err != nil {
			return err
		}

		repository.Save(day, m)
		return nil
	})

	commandHandler.Register(commands.CancelSlotBooking{}, func(c infrastructure.Command, m infrastructure.CommandMetadata) error {
		cmd := c.(commands.CancelSlotBooking)
		day, err := repository.Get(NewDayIDFrom(cmd.DayId))
		if err != nil {
			return err
		}

		err = day.CancelSlotBooking(NewSlotID(cmd.SlotId), cmd.Reason)
		if err != nil {
			return err
		}

		repository.Save(day, m)
		return nil
	})

	commandHandler.Register(commands.ScheduleSlot{}, func(c infrastructure.Command, m infrastructure.CommandMetadata) error {
		cmd := c.(commands.ScheduleSlot)
		day, err := repository.Get(NewDayID(NewDoctorID(cmd.DoctorId), cmd.Start))
		if err != nil {
			return err
		}

		err = day.ScheduleSlot(NewSlotID(cmd.SlotId).Value, cmd.Start, cmd.Duration)
		if err != nil {
			return err
		}

		repository.Save(day, m)
		return nil
	})

	commandHandler.Register(commands.ArchiveDaySchedule{}, func(c infrastructure.Command, m infrastructure.CommandMetadata) error {
		cmd := c.(commands.ArchiveDaySchedule)
		day, err := repository.Get(NewDayIDFrom(cmd.DayId))
		if err != nil {
			return err
		}

		err = day.Archive()
		if err != nil {
			return err
		}

		repository.Save(day, m)
		return nil
	})

	return commandHandler
}
