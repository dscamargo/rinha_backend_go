package controllers

import (
	"errors"
	"fmt"
	"github.com/dscamargo/rinha_backend_go/pessoa"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type PessoaController struct {
	service pessoa.UseCase
}

func NewPessoaController(service pessoa.UseCase) *PessoaController {
	return &PessoaController{service}
}

type CreatePessoaRequest struct {
	Nome       string
	Apelido    string
	Nascimento string
	Stack      []string
}

func (body *CreatePessoaRequest) validate() error {
	if body.Nome == "" {
		return pessoa.ErrInvalidBody
	}

	if len(body.Nome) > 100 {
		return pessoa.ErrInvalidBody
	}

	if body.Apelido == "" {
		return pessoa.ErrInvalidBody
	}

	if len(body.Apelido) > 32 {
		return pessoa.ErrInvalidBody
	}

	if _, err := time.Parse("2006-01-02", body.Nascimento); err != nil {
		return pessoa.ErrInvalidBody
	}

	for _, stack := range body.Stack {
		if len(stack) > 32 {
			return pessoa.ErrInvalidBody
		}
	}

	return nil
}

func (p *PessoaController) Create(c *fiber.Ctx) error {
	var body CreatePessoaRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "bad_request"})
	}

	if err := body.validate(); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	createdPessoa, err := p.service.Create(
		body.Nome,
		body.Apelido,
		body.Nascimento,
		body.Stack,
	)
	if err != nil {
		if errors.Is(err, pessoa.ErrApelidoJaUtilizado) {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Set("Location", fmt.Sprintf("/pessoas/%s", createdPessoa.ID))
	c.Status(http.StatusCreated)
	return c.JSON(mapPessoaOutput(createdPessoa))
}

func (p *PessoaController) List(c *fiber.Ctx) error {
	term := c.Query("t")

	if term == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_param"})
	}

	ps, err := p.service.Search(term)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	formattedPs := make([]PessoaOutput, 0)
	for _, item := range ps {
		formattedPs = append(formattedPs, mapPessoaOutput(&item))
	}

	return c.Status(http.StatusOK).JSON(formattedPs)
}
func (p *PessoaController) Show(c *fiber.Ctx) error {
	id := c.Params("id")
	ps, err := p.service.FindById(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "not_found"})
	}
	return c.Status(http.StatusOK).JSON(mapPessoaOutput(ps))

}
func (p *PessoaController) Count(c *fiber.Ctx) error {
	total, err := p.service.Count()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(http.StatusOK).Send([]byte(fmt.Sprintf("%d", total)))

}
